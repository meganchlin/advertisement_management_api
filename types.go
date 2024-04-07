package main

import "time"

type Gender string

const (
	Male   Gender = "M"
	Female Gender = "F"
)

type Country string

const (
	Taiwan        Country = "TW"
	Japan         Country = "JP"
	United_States Country = "US"
	Korea         Country = "KR"
	Thailand      Country = "TH"
)

type Platform string

const (
	IOS     Platform = "ios"
	Android Platform = "android"
	Web     Platform = "web"
)

// Condition represents data about a record condition.
type Condition struct {
	AgeStart int        `json:"ageStart" bson:"ageStart"`
	AgeEnd   int        `json:"ageEnd" bson:"ageEnd"`
	Gender   []Gender   `json:"gender" bson:"gender"`
	Country  []Country  `json:"country" bson:"country"`
	Platform []Platform `json:"platform" bson:"platform"`
}

// Advertisement represents data about a record advertisement.
type Advertisement struct {
	Title      string      `json:"title" bson:"title"`
	StartAt    time.Time   `json:"startAt" bson:"startAt"`
	EndAt      time.Time   `json:"endAt" bson:"endAt"`
	Conditions []Condition `json:"conditions" bson:"conditions"`
}

// define the sructure of Public API response
type DisplayAds struct {
	Items []AdItem `json:"items" bson:"items"`
}

// AdItem represents data about a record of an ad to be displayed.
type AdItem struct {
	Title string    `json:"title" bson:"title"`
	EndAt time.Time `json:"endAt" bson:"endAt"`
}
