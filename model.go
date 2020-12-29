package swdocs

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type SwDoc struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Created     *time.Time   `json:"created,omitempty"`
	Updated     *time.Time   `json:"updated,omitempty"`
	Description string       `json:"description"`
	Sections    SectionSlice `json:"sections,omitempty"`
}

type SectionSlice []Section

type Section struct {
	Header string    `json:"header"`
	Links  LinkSlice `json:"links"`
}

type LinkSlice []Link

type Link struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// Value - Implementation of valuer for database/sql
func (s SectionSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}
