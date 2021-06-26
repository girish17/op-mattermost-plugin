package util

type ActivityLinks struct {
	Self      Self   `json:"self"`
	Projects  []Self `json:"projects"`
}

type AllowedValues struct {
	Type       string           `json:"type"`
	Id         int              `json:"id"`
	Name       string           `json:"name"`
	Position   int              `json:"position"`
	Default    bool             `json:"default"`
	Links      ActivityLinks    `json:"_links"`
}

type Activity struct {
	Type       string           `json:"type"`
	Name       string           `json:"name"`
	Required   bool             `json:"required"`
	HasDefault bool             `json:"hasDefault"`
	Writable   bool             `json:"writable"`
	Location   string           `json:"location"`
	Embedded   struct{
		AllowedValues []AllowedValues `json:"allowedValues"`
	} `json:"_embedded"`
}

type Dependency struct {

}

type Rtl struct {

}

type TimeEntryOption struct {
	Rtl  Rtl `json:"rtl"`
}

type AllowedValue struct {
	Href string `json:"href"`
}

type LinksTimeEntryWP struct {
	AllowedValues AllowedValue `json:"allowedValues"`
}

type TimeEntryWP struct {
	Type       string           `json:"type"`
	Name       string           `json:"name"`
	Required   bool             `json:"required"`
	HasDefault bool             `json:"hasDefault"`
	Writable   bool             `json:"writable"`
	Location   string           `json:"location"`
	Links      LinksTimeEntryWP  `json:"_links"`
}

type Id struct {
	Type       string          `json:"type"`
	Name       string          `json:"name"`
	Required   bool            `json:"required"`
	HasDefault bool            `json:"hasDefault"`
	Writable   bool            `json:"writable"`
	Options    TimeEntryOption `json:"options"`
}

type SchemaLinks struct {

}

type Schema struct {
	Type         string        `json:"_type"`
	Dependencies []Dependency  `json:"_dependencies"`
	Id           Id            `json:"id"`
	CreatedAt    Id            `json:"createdAt"`
	UpdatedAt    Id            `json:"updatedAt"`
	SpentOn      Id            `json:"spentOn"`
	Hours        Id            `json:"hours"`
	User         Id            `json:"user"`
	WorkPackage  TimeEntryWP   `json:"workPackage"`
	Project      TimeEntryWP   `json:"project"`
	Activity     Activity      `json:"activity"`
	CustomField1 Id            `json:"customField1"`
	Links        SchemaLinks   `json:"_links"`
}

type Link struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

type PayloadLinks struct {
	Project     Link    `json:"project"`
	Activity    Link    `json:"activity"`
	WorkPackage Link    `json:"workPackage"`
}

type Comment struct {
	Format string `json:"format"`
	Raw    string `json:"raw"`
	Html   string `json:"html"`
}

type Payload struct {
	Links        PayloadLinks `json:"_links"`
	Hours        string       `json:"hours"`
	Comment      Comment      `json:"comment"`
	SpentOn      string       `json:"spentOn"`
	CustomField1 Comment      `json:"customField1"`
}

type ValidationError struct {

}

type EmbeddedTimeEntry struct {
    Payload           Payload         `json:"payload"`
    Schema            Schema          `json:"schema"`
    ValidationErrors  ValidationError `json:"validationErrors"`
}

type LinksTimeEntry struct {
	Self     Self `json:"self"`
	Validate Self `json:"validate"`
	Commit   Self `json:"commit"`
}

type TimeEntries struct {
	Type              string            `json:"_type"`
	Embedded          EmbeddedTimeEntry `json:"_embedded"`
	Links             LinksTimeEntry    `json:"_links"`
}

type TimeEntriesBody struct {
	Links struct{
		WorkPackage struct{
			Href string `json:"href"`
		} `json:"workPackage"`
	} `json:"_links"`
}

type UpdateImmediately struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

type Delete struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

type TimeLinks struct {
	Self              Self              `json:"self"`
	UpdateImmediately UpdateImmediately `json:"updateImmediately"`
	Delete            Delete            `json:"delete"`
	Project           Link              `json:"project"`
	WorkPackage       Link              `json:"workPackage"`
	User              Link              `json:"user"`
	Activity          Link              `json:"activity"`
	CustomField1      Link              `json:"CustomField1"`
}

type TimeElement struct {
	Type         string    `json:"_type"`
	Id           int       `json:"id"`
	Comment      Comment   `json:"comment"`
	SpentOn      string    `json:"spentOn"`
	Hours        string    `json:"hours"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
	Links        TimeLinks `json:"_links"`
	CustomField1 float64   `json:"customField1"`
}

type TimeEntryList struct {
	Type              string            `json:"_type"`
	Total             int               `json:"total"`
	Count             int               `json:"count"`
	PageSize          int               `json:"pageSize"`
	Offset            int               `json:"offset"`
	Embedded struct{
		Elements []TimeElement `json:"elements"`
	} `json:"_embedded"`
}
