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
	regexID := regexp.MustCompile(`^[0-9].{6,9}$`)
	isOk := regexID.MatchString(cleanedUserIDStr)
	if !isOk {
		return errors.New(fmt.Sprintf("%s: %s", ErrInvalidUserID, userIDStr))
	}
	return nil
}
