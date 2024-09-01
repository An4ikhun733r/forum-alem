package app

import (
	"bytes"
	"fmt"
	"forum/internal/models"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *Application) RenderError(w http.ResponseWriter, status int, message string) {
	data := &models.TemplateData{
		ErrorCode: status,
		Message:   message,
	}

	app.Render(w, status, "error.tmpl", data)
}

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	// Render the custom error template
	app.RenderError(w, http.StatusInternalServerError, "An unexpected error occurred. Please try again later.")
}

// The clientError helper sends a specific status code and corresponding description
// to the user. We'll use this later in the book to send responses like 400 "Bad
// Request" when there's a problem with the request that the user sent.
func (app *Application) ClientError(w http.ResponseWriter, status int) {
	// Render the custom error template
	app.RenderError(w, status, http.StatusText(status))
}

// For consistency, we'll also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to
// the user.
func (app *Application) NotFound(w http.ResponseWriter) {
	app.RenderError(w, http.StatusNotFound, "Sorry, the page you are looking for does not exist.")
}

func (app *Application) Render(w http.ResponseWriter, status int, page string, data *models.TemplateData) {
	ts, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(w, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
}

func (app *Application) NewTemplateData(r *http.Request) *models.TemplateData {
	return &models.TemplateData{
		IsAuthenticated: false,
		CurrentYear:     time.Now().Year(),
	}
}
