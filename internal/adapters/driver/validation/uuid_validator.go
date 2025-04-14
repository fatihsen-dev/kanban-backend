package validation

import "github.com/go-playground/validator/v10"

type UUIDValidator struct {
	UUID string `json:"uuid" validate:"required,uuid4"`
}

func ValidateUUID(uuid string) error {
	uuidValidator := UUIDValidator{UUID: uuid}
	validate := validator.New()
	err := validate.Struct(uuidValidator)
	if err != nil {
		return err
	}
	return nil
}
