package swdocs

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type swDoc struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	User        string       `json:"user"`
	Created     timeStamp    `json:"created,omitempty"`
	Updated     timeStamp    `json:"updated,omitempty"`
	Description string       `json:"description"`
	Sections    sectionSlice `json:"sections,omitempty"`
}

type swDocsSlice struct {
	SwDocs *[]swDoc
}

type timeStamp time.Time

type sectionSlice []section

type section struct {
	Header      string    `json:"header"`
	Description string    `json:"description"`
	Links       linkSlice `json:"links"`
}

type linkSlice []link

type link struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

// Value - Implementation of valuer for database/sql
func (s sectionSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *sectionSlice) Scan(v interface{}) error {
	var data []byte
	if b, ok := v.([]byte); ok {
		data = b
	} else if s, ok := v.(string); ok {
		data = []byte(s)
	}
	return json.Unmarshal(data, s)
}

func (t *timeStamp) Scan(v interface{}) error {
	// Should be more strictly to check this type.
	vt, err := time.Parse("2006-01-02 15:04:05", v.(string))
	if err != nil {
		return err
	}
	*t = timeStamp(vt)
	return nil
}

func (t timeStamp) ToString() string {
	return time.Time(t).Format("2006-01-02")
}
