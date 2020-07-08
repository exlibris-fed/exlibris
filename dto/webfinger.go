package dto

// A Webfinger is a response to the webfinger endpoint (.well-known/webfinger) to dereference a username.
type Webfinger struct {
	Subject string          `json:"subject"`
	Aliases []string        `json:"aliases"`
	Links   []WebfingerLink `json:"links"`
}

// A WebfingerLink is a structured link in a Webfinger.
type WebfingerLink struct {
	Rel      string `json:"rel"`
	Type     string `json:"type,omitempty"`
	Href     string `json:"href,omitempty"`
	Template string `json:"template,omitempty"`
}
