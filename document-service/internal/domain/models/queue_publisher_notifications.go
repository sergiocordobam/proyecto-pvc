package models

type NotificationMessage struct {
	Event     string                 `json:"event"`
	User      int                    `json:"user"`
	Name      string                 `json:"name"`
	UserEmail string                 `json:"user_email"`
	ExtraData map[string]interface{} `json:"extra_data"`
}
