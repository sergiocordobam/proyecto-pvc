package validator

type ReqValidatorInterface interface {
	ValidateUserID(userID int) error
}
