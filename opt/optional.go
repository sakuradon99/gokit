package opt

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Optional[T any] struct {
	value  T
	exists bool
}

// Empty creates a new Optional without a value.
func Empty[T any]() Optional[T] {
	var empty T
	return Optional[T]{empty, false}
}

// Of creates a new Optional with a value.
func Of[T any](value T) Optional[T] {
	return Optional[T]{value, true}
}

func OfNullable[T any](value *T) Optional[T] {
	if value == nil {
		return Empty[T]()
	}
	return Of[T](*value)
}

// Exists Returns true if the optional value has been set
func (o *Optional[T]) Exists() bool {
	return o.exists
}

func (o *Optional[T]) Get() T {
	return o.value
}

// GetAndExists returns the value and whether it exists.
// It's invalid to use the returned value if the bool is false.
func (o *Optional[T]) GetAndExists() (T, bool) {
	return o.value, o.exists
}

// GetOrElse returns the value if it exists and returns defaultValue otherwise.
func (o *Optional[T]) GetOrElse(defaultValue T) T {
	if !o.exists {
		return defaultValue
	}
	return o.value
}

// MustGet returns the value if it exists and panics otherwise.
func (o *Optional[T]) MustGet() T {
	if !o.exists {
		panic(".MustGet() called on optional Optional value that doesn't exist.")
	}
	return o.value
}

func (o *Optional[T]) GetPointer() *T {
	if !o.exists {
		return nil
	}
	return &o.value
}

// MarshalJSON implements the json.Marshaller interface. Optionals wil
// marshall and unmarshall as a nullable json field. Value type must also
// implement json.Marshaller.
func (o *Optional[Value]) MarshalJSON() ([]byte, error) {
	if !o.exists {
		return json.Marshal(nil)
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface. Optionals will
// marshall and unmarshall as a nullable json field. Value type must also
// implement json.Unmarshaler.
func (o *Optional[Value]) UnmarshalJSON(data []byte) error {
	var v *Value

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v == nil {
		o.exists = false
		return nil
	}

	o.exists = true
	o.value = *v
	return nil
}

// Scan implements the Scanner interface.
func (o *Optional[Value]) Scan(value any) error {
	if value == nil {
		o.exists = false
		return nil
	}

	val, ok := value.(Value)
	if !ok {
		return fmt.Errorf("failed to scan a '%v' into an Optional", value)
	}

	o.exists = true
	o.value = val

	return nil
}

// Value implements the Valuer interface.
func (o *Optional[Value]) Value() (driver.Value, error) {
	if !o.exists {
		return nil, nil
	}
	return o.value, nil
}
