package main

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

func (g Gender) IsValid() bool {
	switch g {
	case Male, Female:
		return true
	default:
		return false
	}
}

func (g *Gender) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	gender := Gender(strings.ToUpper(s))
	if !gender.IsValid() {
		return errors.New("invalid gender value")
	}
	*g = gender
	return nil
}

func (c Country) IsValid() bool {
	switch c {
	case Taiwan, Japan, United_States, Korea, Thailand:
		return true
	default:
		return false
	}
}

func (c *Country) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	country := Country(strings.ToUpper(s))
	if !country.IsValid() {
		return errors.New("invalid gender value")
	}
	*c = country
	return nil
}

func (p Platform) IsValid() bool {
	switch p {
	case IOS, Android, Web:
		return true
	default:
		return false
	}
}

func (p *Platform) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	platform := Platform(strings.ToLower(s))
	if !platform.IsValid() {
		return errors.New("invalid gender value")
	}
	*p = platform
	return nil
}

// ParseTime parses a string in the expected time format.
func ParseTime(s string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	return time.Parse(layout, s)
}

// Customizes the JSON marshalling behavior for the Advertisement struct
func (ad Advertisement) MarshalJSON() ([]byte, error) {
	// Create a map to hold the serialized data
	data := map[string]interface{}{
		"title":      ad.Title,
		"startAt":    ad.StartAt.Format("2006-01-02T15:04:05.000Z"),
		"endAt":      ad.EndAt.Format("2006-01-02T15:04:05.000Z"),
		"conditions": ad.Conditions,
	}

	// Marshal the map to JSON
	return json.Marshal(data)
}
