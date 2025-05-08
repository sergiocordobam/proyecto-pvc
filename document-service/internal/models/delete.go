package models

type DeleteRequest struct {
	UserID    int      `json:"user_id"`
	FileNames []string `json:"files"`
}
