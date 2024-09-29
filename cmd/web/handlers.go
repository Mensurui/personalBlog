package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

type articleCreateForm struct {
	Title   string `form:"name"`
	Content string `form:"content"`
}

type userCreateForm struct {
	Username string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type userLoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	articles, err := app.articleModel.Latest()

	if err != nil {
		http.Error(w, "Unable to retrieve article", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := app.newTemplateData(r)
	data.Articles = articles
	data.IsAuthenticated = app.isAuthenticated(r)
	if !data.IsAuthenticated {
		log.Println("It hasnt been verified")
	}

	app.render(w, http.StatusOK, "home.tmpl", data)

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "about.tmpl", nil)
}

func (app *application) articleView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		log.Println("No ID found")
	}
	article, err := app.articleModel.Get(id)

	if err != nil {
		log.Printf("Error is: %s", err)
	}

	data := app.newTemplateData(r)
	data.Article = article
	app.render(w, http.StatusOK, "article.tmpl", data)
}

func (app *application) write(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "write.tmpl", data)
}

func (app *application) writePost(w http.ResponseWriter, r *http.Request) {
	var form articleCreateForm
	err := r.ParseForm()
	if err != nil {
		log.Printf("The error is: %v", err)
	}
	err = app.form.Decode(&form, r.PostForm)

	if err != nil {
		log.Printf("The error is: %v", err)
	}

	_, err = app.articleModel.Add(form.Title, form.Content)

	if err != nil {
		log.Printf("The erros is: %v", err)
	}

	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) signUpProceed(w http.ResponseWriter, r *http.Request) {
	var form userCreateForm

	err := r.ParseForm()
	if err != nil {
		log.Printf("The error is: %s", err)
	}

	err = app.form.Decode(&form, r.PostForm)
	if err != nil {
		log.Printf("The error is: %v", err)
	}

	err = app.userModel.Create(form.Email, form.Username, form.Password)
	if err != nil {
		log.Printf("The error is: %s", err)
	}
	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.IsAuthenticated = app.isAuthenticated(r)
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := r.ParseForm()
	if err != nil {
		log.Printf("The error is: %s", err)
		http.Error(w, "Unable to process your request", http.StatusBadRequest)
		return
	}

	err = app.form.Decode(&form, r.PostForm)
	if err != nil {
		log.Printf("The error is: %s", err)
		http.Error(w, "Unable to process your request", http.StatusBadRequest)
		return
	}

	// Basic validation
	if form.Email == "" || form.Password == "" {
		log.Printf("Email or password is empty")
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	id, err := app.userModel.Login(form.Email, form.Password)
	if err != nil {
		log.Printf("Login failed for email %s: %s", form.Email, err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		log.Printf("The error is: %s", err)
		http.Error(w, "Unable to process your request", http.StatusInternalServerError)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/write", http.StatusSeeOther)
}
