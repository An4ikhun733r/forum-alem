package models

type TemplateData struct {
	Snippet         *Snippet
	Snippets        []*Snippet
	Comments        *[]Comment
	CurrentYear     int
	Form            any
	IsAuthenticated bool
	User            *User
	ErrorCode       int
	Message         string
}
