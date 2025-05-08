package models

import (
	"fmt"
	"time"
)

const (
	TemporalStatus = "temporal"
	VerfiiedStatus = "verified"
)

type Document struct {
	URL      string   `json:"url,omitempty"`
	Metadata Metadata `json:"metadata"`
}
type Metadata struct {
	Name         string    `json:"name"`
	AbsPath      string    `json:"abs_path"`
	Size         int       `json:"size"`
	OwnerID      int       `json:"owner_id"`
	Type         string    `json:"type"`
	CreationDate time.Time `json:"creation_date"`
	ContentType  string    `json:"content_type"`
	Status       string    `json:"status,omitempty"`
}

func NewDocument(name, documentType, contentType string, size, ownerId int) Document {
	currentDate := time.Now()
	return Document{
		URL: "",
		Metadata: Metadata{
			Name:         name,
			Size:         size,
			OwnerID:      ownerId,
			Type:         documentType,
			CreationDate: currentDate,
			ContentType:  contentType,
			AbsPath:      fmt.Sprintf("%d/%s", ownerId, name),
			Status:       TemporalStatus,
		},
	}
}
