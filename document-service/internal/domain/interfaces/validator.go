package interfaces

type ReqValidatorInterface interface {
	ValidateUserID(userID int) error
}
