package util

type AllowedValues struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type Activity struct {
	Embedded struct{
		AllowedValues []AllowedValues `json:"allowedValues"`
	} `json:"embedded"`
}

type Schema struct {
	Activity Activity `json:"activity"`
}

type TimeEntries struct {
	Type              string            `json:"_type"`
	Links             Links             `json:"_links"`
	Embedded struct{
		Schema   Schema `json:"schema"`
	} `json:"_embedded"`
}
