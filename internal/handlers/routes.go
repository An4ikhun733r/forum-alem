package handlers

import "net/http"

func (h *HandlerApp) Routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", h.Home)
	
	mux.HandleFunc("/snippet/view/", h.SnippetView)
	mux.HandleFunc("/snippet/create", h.RequireAuth(h.SnippetCreate))
	mux.HandleFunc("/snippet/create/post", h.RequireAuth(h.SnippetCreatePost))

	mux.HandleFunc("/user/signup", h.UserSignup)
	mux.HandleFunc("/user/signup/post", h.UserSignupPost)
	mux.HandleFunc("/user/login", h.UserLogin)
	mux.HandleFunc("/user/login/post", h.UserLoginPost)
	mux.HandleFunc("/user/logout", h.UserLogout)

	mux.HandleFunc("/snippet/like", h.RequireAuth(h.LikePost))
	mux.HandleFunc("/snippet/dislike", h.RequireAuth(h.DislikePost))
	mux.HandleFunc("/snippet/comment", h.RequireAuth(h.AddComment))

	// Wrap the existing chain with the recoverPanic middleware.
	return h.recoverPanic(h.logRequest(secureHeaders(mux)))
}
