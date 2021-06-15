package util

type Self struct{
	Href string `json:"href"`
	Title string `json:"title"`
}

type Links struct{
	Self                       Self              `json:"self"`
	CreateWorkPackage          CreateWorkPackage `json:"createWorkPackage"`
	CreateWorkPackageImmediate CreateWorkPackage `json:"createWorkPackageImmediate"`
	Categories                 Categories        `json:"categories"`
	Versions                   Categories        `json:"versions"`
	Projects                   Categories        `json:"projects"`
	Status                     Self              `json:"status"`
}

type CreateWorkPackage struct {
	Href   string `json:"href"`
	Method string `json:"method"`
}

type Categories struct {
	Href string `json:"href"`
}

type StatusExplanation struct {
	Format string `json:"format"`
	Raw    string `json:"raw"`
	Html   string `json:"html"`
}

type Element struct {
	Type              string            `json:"_type"`
	Links             Links             `json:"_links"`
	Id                int               `json:"id"`
	Identifier        string            `json:"identifier"`
	Name              string            `json:"name"`
	Active            string            `json:"active"`
	StatusExplanation StatusExplanation `json:"statusExplanation"`
	Public            bool              `json:"public"`
	Description       StatusExplanation `json:"description"`
	CreatedAt         string            `json:"createdAt"`
	UpdatedAt         string            `json:"updatedAt"`
	ProjectType       string            `json:"type"`
	Subject           string            `json:"subject"`
}

type Projects struct {
	Links    Links  `json:"_links"`
	Type     string `json:"_type"`
	Total    int64  `json:"total"`
	Count    int64  `json:"count"`
	Embedded struct{
		Elements []Element `json:"elements"`
	} `json:"_embedded"`
}
