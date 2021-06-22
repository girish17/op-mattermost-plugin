package util

type WorkPackages struct {
	Type              string            `json:"_type"`
	Links             Links             `json:"_links"`
	Total             int               `json:"total"`
	Count             int               `json:"count"`
	Embedded struct{
		Elements []Element `json:"elements"`
	} `json:"_embedded"`
}
