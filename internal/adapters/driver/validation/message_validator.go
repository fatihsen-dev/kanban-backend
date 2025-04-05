package validation

func ValidateMessage(message interface{}) error {
	err := Validate(message)
	if err != nil {
		return err
	}
	return nil
}
