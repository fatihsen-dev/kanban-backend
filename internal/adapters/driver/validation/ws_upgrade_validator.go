package validation

type UpgradeRequest struct {
	GroupID string `json:"group_id" validate:"required,min=1"`
}

func ValidateWSUpgradeRequest(request *UpgradeRequest) error {
	err := Validate(request)
	if err != nil {
		return err
	}
	return nil
}
