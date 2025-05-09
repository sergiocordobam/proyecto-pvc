package models

type RegisterDocumentsMessage struct {
	CitizenId int      `json:"citizenId"`
	Documents []string `json:"documents"`
}
type DeleteDocumentsMessage struct {
	CitizenId int `json:"citizenId"`
}
