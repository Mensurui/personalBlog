package main

import (
	"github.com/Mensurui/personalBlog.git/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	wrapper := alice.New(app.sessionManager.LoadAndSave)
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/about", app.about)
	router.HandlerFunc(http.MethodGet, "/article/view/:id", app.articleView)
	router.HandlerFunc(http.MethodGet, "/signup", app.signUp)
	router.HandlerFunc(http.MethodPost, "/signup/create", app.signUpProceed)
	router.Handler(http.MethodGet, "/login", wrapper.ThenFunc(app.login))
	router.Handler(http.MethodPost, "/login/authenticate", wrapper.ThenFunc(app.loginPost))

	router.Handler(http.MethodGet, "/write", wrapper.ThenFunc(app.write))
	router.Handler(http.MethodPost, "/writePost", wrapper.ThenFunc(app.writePost))
	return router
}
