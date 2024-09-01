package handlers

import (
	"errors" // New import
	"fmt"
	"forum/internal/app"
	"forum/internal/models"
	"forum/internal/validator" //"html/template"
	"log"
	"net/http"
	"strconv" // New import
	// New import
)

func (h *HandlerApp) Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		h.NotFound(w)
		return
	}

	// Extract and validate query parameters
	query := r.URL.Query()
	selectedTags := query["tag"]
	filter := query.Get("filter")

	// Define allowed tags
	allowedTags := map[string]struct{}{
		"Memes": {},
		"Life":  {},
		"Games": {},
	}

	// Validate 'tag' parameters
	for _, tag := range selectedTags {
		if _, valid := allowedTags[tag]; !valid {
			h.ClientError(w, http.StatusBadRequest)
			return
		}
	}

	// Validate 'filter' parameter
	if filter != "" && filter != "liked" && filter != "myPosts" {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	// Redirect to login if the user is not authenticated and tries to filter by 'liked' or 'myPosts'
	if (filter == "liked" || filter == "myPosts") && !h.IsAuthenticated(r) {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// Initialize userID
	var userID int

	// Retrieve userID if the user is authenticated
	if h.IsAuthenticated(r) {
		user, err := h.service.GetUser(r) // Get the *models.User object
		if err != nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		userID = user.ID // Extract the user ID from the User object
	}

	// Retrieve snippets with applied filters
	snippets, err := h.service.Latest(selectedTags, filter, userID)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Prepare data for rendering
	data := h.NewTemplateData(r)
	data.IsAuthenticated = h.IsAuthenticated(r)
	data.Snippets = snippets

	h.Render(w, http.StatusOK, "home.tmpl", data)
}

func (h *HandlerApp) SnippetView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Path[len("/snippet/view/"):]

	// Check if the ID part is empty
	if idStr == "" {
		h.NotFound(w) // No ID provided, so return 404 Not Found
		return
	}

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 { // ID is invalid if <= 0 or conversion failed
		h.ClientError(w, http.StatusBadRequest) // Return 400 Bad Request for invalid IDs
		return
	}

	snippet, err := h.service.GetSnippet(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFound(w)
		} else {
			h.ServerError(w, err)
		}
		return
	}
	comments, err := h.service.GetCommentByPostId(id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	data := h.NewTemplateData(r)
	data.IsAuthenticated = h.IsAuthenticated(r)
	data.Snippet = snippet
	data.Comments = &comments

	h.Render(w, http.StatusOK, "view.tmpl", data)
}

func (h *HandlerApp) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	data := h.NewTemplateData(r)
	data.IsAuthenticated = h.IsAuthenticated(r)
	data.Form = models.SnippetCreateForm{}

	h.Render(w, http.StatusOK, "create.tmpl", data)
}

func (h *HandlerApp) SnippetCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form := models.SnippetCreateForm{
		Title:    r.PostForm.Get("title"),
		Content:  r.PostForm.Get("content"),
		Category: r.PostForm["category"],
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	if !form.Valid() {
		data := h.NewTemplateData(r)
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	cookies := app.GetSessionCookie("session_id", r)

	id, err := h.service.InsertSnippet(cookies.Value, form.Title, form.Content, form.Category)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (h *HandlerApp) UserSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	data := h.NewTemplateData(r)
	data.Form = models.UserSignupForm{}
	h.Render(w, http.StatusOK, "signup.tmpl", data)
}

func (h *HandlerApp) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form := models.UserSignupForm{
		Name:     r.PostForm.Get("name"),
		Password: r.PostForm.Get("password"),
		Email:    r.PostForm.Get("email"),
	}

	form.CheckField(validator.NotBlank(form.Name), "username", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, models.EmailRX), "email", "This field must be a valid email address")

	if !form.Valid() {
		data := h.NewTemplateData(r)
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	_, err = h.service.InsertUser(form.Name, form.Password, form.Email)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("generic", "Such Email already registred")
			data := h.NewTemplateData(r)
			data.Form = form
			h.Render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			h.ServerError(w, err)
		}
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// User login handlers
func (h *HandlerApp) UserLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	data := h.NewTemplateData(r)
	data.Form = models.UserLoginForm{}
	h.Render(w, http.StatusOK, "login.tmpl", data)
}

func (h *HandlerApp) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	form := models.UserLoginForm{
		Name:     r.PostForm.Get("name"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Name), "username", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := h.NewTemplateData(r)
		data.Form = form
		h.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	session, _, err := h.service.Authenticate(form.Name, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("generic", "Username or password is incorrect")
			data := h.NewTemplateData(r)
			data.Form = form
			h.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			fmt.Println(err)
			h.ServerError(w, err)
		}
		return
	}

	data := h.NewTemplateData(r)
	data.IsAuthenticated = true
	app.SetSessionCookie("session_id", w, session.Token, session.ExpTime)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerApp) UserLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	c := app.GetSessionCookie("session_id", r)
	if c != nil {
		h.service.DeleteSession(c.Value)
		app.ExpireSessionCookie("session_id", w)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerApp) LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID < 1 {
		log.Println(err)
		h.NotFound(w)
		return
	}

	userID, err := h.service.GetUser(r)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFound(w)
			return
		}
		h.ServerError(w, err)
		return
	}

	err = h.service.LikePost(userID.ID, postID)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerApp) DislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.FormValue("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID < 1 {
		log.Println(err)
		h.NotFound(w)
		return
	}

	userID, err := h.service.GetUser(r)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			h.NotFound(w)
			return
		}
		h.ServerError(w, err)
		return
	}

	err = h.service.DislikePost(userID.ID, postID)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerApp) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodPost {
		postIDStr := r.FormValue("PostId")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil || postID < 1 {
			h.NotFound(w)
			return
		}
		userID, err := h.service.GetUser(r)
		content := r.FormValue("Content")
		if err != nil {
			h.ServerError(w, err)
			return
		}
		err = h.service.AddComment(postID, userID.ID, content)
		if err != nil {
			h.ClientError(w, http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", postID), http.StatusSeeOther)
	}
}
