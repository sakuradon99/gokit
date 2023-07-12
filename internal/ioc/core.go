package ioc

import (
	"fmt"
	"github.com/sakuradon99/gokit/internal/db"
	"github.com/spf13/viper"
	"reflect"
	"unsafe"
)

type Container interface {
	Register(object any, opts ...RegisterOption) error
	GetObject(name string, t any) (any, error)
}

type ContainerImpl struct {
	objectPool    *ObjectPool
	interfacePool *InterfacePool
	loaded        bool
}

func NewContainerImpl() *ContainerImpl {
	return &ContainerImpl{
		objectPool:    NewObjectPool(),
		interfacePool: NewInterfacePool(),
	}
}

func (c *ContainerImpl) Register(object any, opts ...RegisterOption) error {
	var options RegisterOptions
	for _, opt := range opts {
		opt(&options)
	}

	ot := reflect.TypeOf(object)

	if ot.Kind() != reflect.Ptr && ot.Kind() != reflect.Func {
		return fmt.Errorf("unsupported register type %s", ot.Kind())
	}

	var objectID string
	var obj *Object

	if ot.Kind() == reflect.Func {
		if ot.NumOut() > 2 || ot.NumOut() == 2 && ot.Out(1).Name() != "error" {
			return fmt.Errorf("unsupported function")
		}
		ret := ot.Out(0).Elem()
		objectID = genObjectID(ret.PkgPath(), ret.Name(), options.Name)
		obj = NewObjectFromFunc(objectID, options.Name, object)
	} else {
		ot = ot.Elem()
		objectID = genObjectID(ot.PkgPath(), ot.Name(), options.Name)
		obj = NewObject(objectID, options.Name, object)
	}

	err := c.objectPool.Add(obj)
	if err != nil {
		return err
	}

	for _, implementInterface := range options.ImplementInterfaces {
		it := reflect.TypeOf(implementInterface).Elem()
		infID := genInterfaceID(it.PkgPath(), it.Name())
		c.interfacePool.Add(Interface{
			id: infID,
		})
		err = c.interfacePool.BindImpl(infID, objectID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ContainerImpl) GetObject(name string, t any) (any, error) {
	err := c.load()
	if err != nil {
		return nil, err
	}

	tt := reflect.TypeOf(t).Elem()
	id := genObjectID(tt.PkgPath(), tt.Name(), name)

	object, ok := c.objectPool.Get(id)
	if !ok {
		return nil, fmt.Errorf("object %s not found", id)
	}

	return object.obj, nil
}

func (c *ContainerImpl) load() error {
	if c.loaded {
		return nil
	}

	err := c.registerDB()
	if err != nil {
		return err
	}

	err = c.readConfig()
	if err != nil {
		return err
	}

	requiredObjects := c.objectPool.List()

	for _, object := range requiredObjects {
		if object.injected {
			continue
		}
		if object.optional {
			continue
		}
		err := c.inject(object)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ContainerImpl) inject(object *Object) error {
	if object.isFunc {
		return c.injectFunc(object)
	} else {
		return c.injectStruct(object)
	}
}

func (c *ContainerImpl) injectFunc(object *Object) error {
	ot := object.GetReflectType()
	ov := object.GetReflectValue()

	incomes := make([]reflect.Value, ot.NumIn())

	for i := 0; i < ot.NumIn(); i++ {
		arg := ot.In(i)
		id := genObjectID(arg.PkgPath(), arg.Name(), "")
		dependency, ok := c.objectPool.Get(id)
		if !ok {
			return fmt.Errorf("missing dependency object %s", id)
		}

		if !dependency.injected {
			err := c.inject(dependency)
			if err != nil {
				return err
			}
		}

		incomes[i] = dependency.GetReflectValue()
	}

	outcomes := ov.Call(incomes)
	if len(outcomes) == 2 && !outcomes[1].IsNil() {
		return outcomes[1].Interface().(error)
	}

	object.SetInjected(true)
	object.SetFuncRet(outcomes[0])

	return nil
}

func (c *ContainerImpl) injectStruct(object *Object) error {
	ot := object.GetReflectType()
	ov := object.GetReflectValue()

	for i := 0; i < ot.NumField(); i++ {
		field := ot.Field(i)
		if requiredObjectName, ok := field.Tag.Lookup("inject"); ok {
			switch field.Type.Kind() {
			case reflect.Interface:
				id := genInterfaceID(field.Type.PkgPath(), field.Type.Name())
				implObjectIDs, err := c.interfacePool.GetImplObjectIDs(id)
				if err != nil {
					return err
				}

				if len(implObjectIDs) == 0 {
					return fmt.Errorf("missing implementation for interface %s", id)
				}

				dependency, ok := c.objectPool.Get(implObjectIDs[0])
				if !ok {
					return fmt.Errorf("missing dependency object %s", implObjectIDs[0])
				}

				if !dependency.injected {
					err = c.inject(dependency)
					if err != nil {
						return err
					}
				}

				assignPrivateField(ov.Field(i), dependency.GetValue())
			case reflect.Struct, reflect.Ptr:
				ft := field.Type
				if field.Type.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}
				id := genObjectID(ft.PkgPath(), ft.Name(), requiredObjectName)
				dependency, ok := c.objectPool.Get(id)
				if !ok {
					return fmt.Errorf("missing dependency object %s", id)
				}
				if !dependency.injected {
					err := c.inject(dependency)
					if err != nil {
						return err
					}
				}

				assignPrivateField(ov.Field(i), dependency.GetValue())
			}
		} else if valueName, ok := field.Tag.Lookup("value"); ok {
			value := viper.Get(valueName)
			if value == nil {
				return fmt.Errorf("missing value %s", valueName)
			}
			assignPrivateField(ov.Field(i), value)
		}
	}

	object.SetInjected(true)
	return nil
}

func (c *ContainerImpl) registerDB() error {
	err := c.Register(new(db.Config), Optional())
	if err != nil {
		return err
	}
	err = c.Register(db.InitGorm, Optional())
	if err != nil {
		return err
	}
	err = c.Register(new(db.ManagerImpl), Implement(new(db.Manager)), Optional())
	if err != nil {
		return err
	}
	return nil
}

func (c *ContainerImpl) readConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func assignPrivateField(field reflect.Value, val any) {
	field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
	if v, ok := val.(reflect.Value); ok {
		field.Set(v)
		return
	}
	field.Set(reflect.ValueOf(val))
}
