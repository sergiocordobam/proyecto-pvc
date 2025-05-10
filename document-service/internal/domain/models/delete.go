package models

type DeleteRequest struct {
	UserID   int    `json:"user_id"`
	FileName string `json:"file"`
	Email    string `json:"email"`
}
