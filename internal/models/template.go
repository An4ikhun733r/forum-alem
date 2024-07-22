package models

type TemplateData struct {
	Snippet     *Snippet
	Snippets    []*Snippet
	CurrentYear int
	Form            any
	IsAuthenticated bool
	User            *User
}