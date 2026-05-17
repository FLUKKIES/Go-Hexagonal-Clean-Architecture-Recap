package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Initialize the validator instance automatically when the package is loaded
func init() {
	validate = validator.New()
	
	// ใช้ชื่อ json tag เป็นชื่อ field ในข้อความ error
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ErrorResponse represents the detailed validation error for a specific field
type ErrorResponse struct {
	FailedField string `json:"field"`   // ชื่อ field ที่มีปัญหา
	Tag         string `json:"tag"`     // ประเภทของ tag (เช่น required, min, max)
	Value       string `json:"value"`   // ค่าที่ส่งมา
	Message     string `json:"message"` // ข้อความ error ที่อ่านง่ายขึ้น
}

// ValidateStruct validates a struct and returns detailed errors
func ValidateStruct(data interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()

			// แปลงข้อความ error ให้อ่านง่ายขึ้น (Custom Error Message)
			switch element.Tag {
			case "required":
				element.Message = fmt.Sprintf("Field '%s' is required", element.FailedField)
			case "email":
				element.Message = fmt.Sprintf("Field '%s' must be a valid email address", element.FailedField)
			case "min":
				element.Message = fmt.Sprintf("Field '%s' must be at least %s characters", element.FailedField, element.Value)
			case "max":
				element.Message = fmt.Sprintf("Field '%s' must be at most %s characters", element.FailedField, element.Value)
			case "len":
				element.Message = fmt.Sprintf("Field '%s' must be exactly %s characters", element.FailedField, element.Value)
			case "numeric":
				element.Message = fmt.Sprintf("Field '%s' must contain only numbers", element.FailedField)
			case "e164":
				element.Message = fmt.Sprintf("Field '%s' must be a valid E.164 phone number (e.g. +66812345678)", element.FailedField)
			default:
				element.Message = fmt.Sprintf("Field '%s' failed validation on the '%s' tag", element.FailedField, element.Tag)
			}

			errors = append(errors, &element)
		}
	}
	return errors
}

// Validate ตรวจสอบ struct และคืนค่า error แบบ key-value สำหรับ response
func Validate(data interface{}) map[string]string {
	errs := ValidateStruct(data)
	if len(errs) > 0 {
		errorMap := make(map[string]string)
		for _, e := range errs {
			errorMap[e.FailedField] = e.Message
		}
		return errorMap
	}
	return nil
}
