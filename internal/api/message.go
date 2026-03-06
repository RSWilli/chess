package api

const EventMarkup = "markup"

type Markup struct {
	Selector string `json:"selector"`
	Markup   string `json:"markup"`
}
