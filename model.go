package swdocs

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type SwDoc struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Created     TimeStamp    `json:"created,omitempty"`
	Updated     TimeStamp    `json:"updated,omitempty"`
	Description string       `json:"description"`
	Sections    SectionSlice `json:"sections,omitempty"`
}

type TimeStamp time.Time

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

func (t *TimeStamp) Scan(v interface{}) error {
	// Should be more strictly to check this type.
	vt, err := time.Parse("2006-01-02 15:04:05", v.(string))
	if err != nil {
		return err
	}
	*t = TimeStamp(vt)
	return nil
}

func (t TimeStamp) ToString() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}
