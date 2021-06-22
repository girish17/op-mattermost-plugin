package util

type Option struct {
	Text string `json:"text"`
	Value string `json:"value"`
}


type Context struct {
	Action string `json:"action"`
}

type Integration struct {
	Url string `json:"url"`
	Context Context `json:"context"`
}

type Action struct {
	Name string `json:"name"`
	Integration Integration`json:"integration"`
	Type string `json:"type"`
	Options []Option `json:"options"`
}

type Attachment struct {
	Actions []Action `json:"actions"`
}

type OptAttachments struct {
	Attachments []Attachment `json:"attachments"`
}
