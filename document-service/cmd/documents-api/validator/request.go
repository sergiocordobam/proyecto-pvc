package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	ErrInvalidUserID = "invalid user ID"
)

type ReqValidator struct {
}

func NewReqValidator() *ReqValidator {
	return &ReqValidator{}
}

func (r ReqValidator) ValidateUserID(userID int) error {
	if userID <= 0 {
		return errors.New(fmt.Sprintf("%s: %d", ErrInvalidUserID, userID))
	}
	userIDStr := strconv.Itoa(userID)
	cleanedUserIDStr := strings.TrimSpace(userIDStr)
	regexiD := regexp.MustCompile(`^[0-9].{6,9}$`)
	isOk := regexiD.MatchString(cleanedUserIDStr)
	if !isOk {
		return errors.New(fmt.Sprintf("%s: %s", ErrInvalidUserID, userIDStr))
	}
	return nil
}
