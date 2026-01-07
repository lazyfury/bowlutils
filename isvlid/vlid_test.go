package isvlid

import (
	"testing"
)

func TestRequired(t *testing.T) {
	tests := []struct {
		name    string
		field   any
		wantErr bool
	}{
		{"zero string", "", true},
		{"zero int", 0, true},
		{"zero float", 0.0, true},
		{"nil pointer", (*int)(nil), true},
		{"non-zero string", "value", false},
		{"non-zero int", 1, false},
		{"non-zero float", 1.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Required()
			err := cond(nil, tt.field, "testField")
			if (err != nil) != tt.wantErr {
				t.Errorf("Required() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsEnum(t *testing.T) {
	enum := []string{"red", "green", "blue"}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid value", "red", false},
		{"valid value 2", "green", false},
		{"invalid value", "yellow", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := IsEnum(enum)
			err := cond(nil, tt.value, "color")
			if (err != nil) != tt.wantErr {
				t.Errorf("IsEnum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsEnum_Int(t *testing.T) {
	enum := []int{1, 2, 3}

	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{
		{"valid value", 1, false},
		{"valid value 2", 2, false},
		{"invalid value", 4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := IsEnum(enum)
			err := cond(nil, tt.value, "number")
			if (err != nil) != tt.wantErr {
				t.Errorf("IsEnum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsOneOf(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"valid value", "a", false},
		{"valid value 2", "b", false},
		{"invalid value", "c", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := IsOneOf("a", "b")
			err := cond(nil, tt.value, "letter")
			if (err != nil) != tt.wantErr {
				t.Errorf("IsOneOf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		name       string
		phone      string
		allowEmpty bool
		wantErr    bool
	}{
		{"valid phone", "13812345678", false, false},
		{"valid phone 2", "15912345678", false, false},
		{"invalid length", "1381234567", false, true},
		{"invalid prefix", "02812345678", false, true},
		{"empty with allow", "", true, false},
		{"empty without allow", "", false, true},
		{"too short", "1381234567", false, true},
		{"too long", "138123456789", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := IsValidPhone(tt.phone, tt.allowEmpty)
			err := cond(nil, tt.phone, "phone")
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidPhone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		allowEmpty bool
		wantErr    bool
	}{
		{"valid email", "test@example.com", false, false},
		{"valid email 2", "user.name+tag@example.co.uk", false, false},
		{"invalid email", "invalid-email", false, true},
		{"invalid email 2", "@example.com", false, true},
		{"invalid email 3", "test@", false, true},
		{"empty with allow", "", true, false},
		{"empty without allow", "", false, true},
		{"missing @", "testexample.com", false, true},
		{"missing domain", "test@", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := IsValidEmail(tt.email, tt.allowEmpty)
			err := cond(nil, tt.email, "email")
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	type TestStruct struct {
		Name  string
		Count int
	}

	tests := []struct {
		name         string
		field        any
		defaultValue any
		shouldSet    bool
	}{
		{"zero string", "", "default", true},
		{"non-zero string", "value", "default", false},
		{"zero int", 0, 10, true},
		{"non-zero int", 5, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var target TestStruct
			if tt.name == "zero string" || tt.name == "non-zero string" {
				target.Name = tt.field.(string)
				cond := Default(tt.defaultValue.(string))
				err := cond(&target, tt.field, "Name")
				if err != nil {
					t.Errorf("Default() error = %v", err)
				}
				if tt.shouldSet && target.Name != tt.defaultValue {
					t.Errorf("Default() should set value, got %v", target.Name)
				}
				if !tt.shouldSet && target.Name != tt.field {
					t.Errorf("Default() should not set value, got %v", target.Name)
				}
			} else {
				target.Count = tt.field.(int)
				cond := Default(tt.defaultValue.(int))
				err := cond(&target, tt.field, "Count")
				if err != nil {
					t.Errorf("Default() error = %v", err)
				}
				if tt.shouldSet && target.Count != tt.defaultValue {
					t.Errorf("Default() should set value, got %v", target.Count)
				}
				if !tt.shouldSet && target.Count != tt.field {
					t.Errorf("Default() should not set value, got %v", target.Count)
				}
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		min     int
		wantErr bool
	}{
		{"valid value", 10, 5, false},
		{"equal to min", 5, 5, false},
		{"less than min", 3, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Min(tt.min)
			err := cond(nil, tt.value, "value")
			if (err != nil) != tt.wantErr {
				t.Errorf("Min() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		max     int
		wantErr bool
	}{
		{"valid value", 5, 10, false},
		{"equal to max", 10, 10, false},
		{"greater than max", 15, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cond := Max(tt.max)
			err := cond(nil, tt.value, "value")
			if (err != nil) != tt.wantErr {
				t.Errorf("Max() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_Validate(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
		Age   int    `validate:"min=18"`
	}

	tests := []struct {
		name    string
		value   *TestStruct
		wantErr bool
	}{
		{"valid struct", &TestStruct{
			Name:  "John",
			Email: "john@example.com",
			Age:   25,
		}, false},
		{"missing name", &TestStruct{
			Email: "john@example.com",
			Age:   25,
		}, true},
		{"invalid email", &TestStruct{
			Name:  "John",
			Email: "invalid-email",
			Age:   25,
		}, true},
		{"age too young", &TestStruct{
			Name:  "John",
			Email: "john@example.com",
			Age:   15,
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator(tt.value)
			err := v.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidator_WithCondition(t *testing.T) {
	type TestStruct struct {
		Status string
	}

	target := &TestStruct{Status: "pending"}

	v := NewValidator(target,
		WithCondition("Status", IsEnum([]string{"pending", "active", "completed"})),
	)

	err := v.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v", err)
	}

	// 测试无效值
	target.Status = "invalid"
	err = v.Validate()
	if err == nil {
		t.Error("Validate() should return error for invalid enum value")
	}
}

func TestValidator_WithUseV10(t *testing.T) {
	type TestStruct struct {
		Name string
	}

	target := &TestStruct{Name: "test"}

	// 测试不使用 v10 validator
	v := NewValidator(target, WithUseV10(false))
	err := v.Validate()
	if err != nil {
		t.Errorf("Validate() error = %v", err)
	}
}
