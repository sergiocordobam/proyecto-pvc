package models

import (
	"time"
)

type Document struct {
	URL      string   `json:"url"`
	Metadata Metadata `json:"metadata"`
}
type Metadata struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	OwnerID      int64     `json:"owner_id"`
	Type         string    `json:"type"`
	CreationDate time.Time `json:"creation_date"`
}

func NewDocument(name, documentType, url string, size, ownerId int64) Document {
	currentDate := time.Now()
	return Document{
		URL: url,
		Metadata: Metadata{
			Name:         name,
			Size:         size,
			OwnerID:      ownerId,
			Type:         documentType,
			CreationDate: currentDate,
		},
	}
}
