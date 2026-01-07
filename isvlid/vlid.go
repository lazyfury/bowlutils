package isvlid

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/lazyfury/bowlutils/utils"
)

type Condition func(target any, field any, fieldName string) error

type Validator struct {
	Value      any
	Conditions map[string][]Condition
	UseV10     bool
}

type ValidatorOption func(v *Validator)

func NewValidator(value any, opts ...ValidatorOption) *Validator {
	v := &Validator{Value: value, UseV10: true}
	for _, opt := range opts {
		opt(v)
	}
	return v
}
func WithUseV10(useV10 bool) ValidatorOption {
	return func(v *Validator) {
		v.UseV10 = useV10
	}
}

func WithCondition(field string, conds ...Condition) ValidatorOption {
	return func(v *Validator) {
		if v.Conditions == nil {
			v.Conditions = make(map[string][]Condition)
		}
		v.Conditions[field] = append(v.Conditions[field], conds...)
	}
}

func (v *Validator) Validate() error {
	if v.UseV10 {
		err := validator.New().Struct(v.Value)
		if err != nil {
			return err
		}
	}
	// logger.Debug("Validate", "Conditions", v.Conditions)
	if len(v.Conditions) == 0 {
		return nil
	}
	// logger.Debug("Validate", "Value", v.Value)
	if reflect.ValueOf(v.Value).Kind() != reflect.Ptr {
		return fmt.Errorf("value must be a pointer")
	}

	value := reflect.ValueOf(v.Value).Elem()
	for field, conds := range v.Conditions {
		// logger.Debug("Validate", "field", field, "conds", conds)
		val := value.FieldByName(field)
		if !val.IsValid() {
			return fmt.Errorf("field %s is not found", field)
		}
		for _, cond := range conds {
			if err := cond(v.Value, val.Interface(), field); err != nil {
				return fmt.Errorf("field %s: %w", field, err)
			}
		}
	}
	return nil
}

func Required() Condition {
	return func(target any, field any, fieldName string) error {
		if utils.IsZero(field) {
			return fmt.Errorf("value is required")
		}
		return nil
	}
}

func IsEnum[T comparable](enum []T) Condition {
	return func(target any, field any, fieldName string) error {
		for _, e := range enum {
			if e == field.(T) {
				return nil
			}
		}
		return fmt.Errorf("value %v is not in enum %v", field.(T), enum)
	}
}

func IsOneOf[T comparable](values ...T) Condition {
	return func(target any, field any, fieldName string) error {
		for _, e := range values {
			if e == field.(T) {
				return nil
			}
		}
		return fmt.Errorf("value %v is not in values %v", field.(T), values)
	}
}

func IsValidPhone(phone string, allowEmpty bool) Condition {
	return func(target any, field any, fieldName string) error {
		if allowEmpty && len(phone) == 0 {
			return nil
		}
		if len(phone) != 11 {
			return fmt.Errorf("phone number must be 11 digits")
		}
		regx := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !regx.MatchString(phone) {
			return fmt.Errorf("phone number is invalid")
		}
		return nil
	}
}

func IsValidEmail(email string, allowEmpty bool) Condition {
	return func(target any, field any, fieldName string) error {
		if allowEmpty && len(email) == 0 {
			return nil
		}
		if len(email) == 0 {
			return fmt.Errorf("email is required")
		}
		regx := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !regx.MatchString(email) {
			return fmt.Errorf("email is invalid")
		}
		return nil
	}
}

func Default[T any](defaultValue T) Condition {
	return func(target any, field any, fieldName string) error {
		if utils.IsZero(field) {
			rv := reflect.ValueOf(target).Elem()
			fv := rv.FieldByName(fieldName)
			if fv.IsValid() {
				fv.Set(reflect.ValueOf(defaultValue))
			}
			return nil
		}
		return nil
	}
}

func Min[T int | int64 | int32 | float32 | float64](min T) Condition {
	return func(target any, field any, fieldName string) error {
		if field.(T) < min {
			return fmt.Errorf("value %v is less than min %v", field.(T), min)
		}
		return nil
	}
}

func Max[T int | int64 | int32 | float32 | float64](max T) Condition {
	return func(target any, field any, fieldName string) error {
		if field.(T) > max {
			return fmt.Errorf("value %v is greater than max %v", field.(T), max)
		}
		return nil
	}
}
