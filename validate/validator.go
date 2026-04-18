package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	engine *validator.Validate
}

func New() *Validator {
	engine := validator.New()
	engine.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("json")
		if tag == "" {
			return field.Name
		}
		name := strings.Split(tag, ",")[0]
		if name == "-" || name == "" {
			return field.Name
		}
		return name
	})
	return &Validator{engine: engine}
}

func (v *Validator) Struct(value any) error {
	if err := v.engine.Struct(value); err != nil {
		return fmt.Errorf("validate struct: %w", err)
	}
	return nil
}

func (v *Validator) Engine() *validator.Validate {
	return v.engine
}
