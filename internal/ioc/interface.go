package ioc

import "fmt"

type Interface struct {
	id            string
	implObjectIDs []string
}

type InterfacePool struct {
	interfaces map[string]Interface
}

func NewInterfacePool() *InterfacePool {
	return &InterfacePool{
		interfaces: make(map[string]Interface),
	}
}

func (p *InterfacePool) Add(inf Interface) {
	if _, ok := p.interfaces[inf.id]; ok {
		return
	}
	p.interfaces[inf.id] = inf
}

func (p *InterfacePool) BindImpl(interfaceID, objectID string) error {
	inf, ok := p.interfaces[interfaceID]
	if !ok {
		return fmt.Errorf("interface with id %s not found", interfaceID)
	}
	inf.implObjectIDs = append(inf.implObjectIDs, objectID)
	p.interfaces[interfaceID] = inf
	return nil
}

func (p *InterfacePool) GetImplObjectIDs(interfaceID string) ([]string, error) {
	inf, ok := p.interfaces[interfaceID]
	if !ok {
		return nil, fmt.Errorf("interface with id %s not found", interfaceID)
	}

	return inf.implObjectIDs, nil
}

func genInterfaceID(pkgPath string, name string) string {
	return fmt.Sprintf("%s.%s", pkgPath, name)
}
