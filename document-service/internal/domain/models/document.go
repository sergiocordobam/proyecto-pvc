package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	TemporalStatus = "temporal"
	VerifiedStatus = "verified"
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

func NewMetadata(name, documentType, contentType string, size, ownerId int) Metadata {
	currentDate := time.Now()
	newName := strings.Replace(name, strconv.Itoa(ownerId)+"/", "", -1)
	return Metadata{
		Name:         newName,
		Size:         size,
		OwnerID:      ownerId,
		Type:         documentType,
		CreationDate: currentDate,
		ContentType:  contentType,
		AbsPath:      fmt.Sprintf("%d/%s", ownerId, newName),
		Status:       TemporalStatus,
	}
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
