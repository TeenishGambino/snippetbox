package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.abiral.net/internal/models"
	"snippetbox.abiral.net/internal/validator"
)

type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// Tells the decoder to ignore the field//
	validator.Validator `form:"-`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we don't need the following
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context. So we use the following

	params := httprouter.ParamsFromContext(r.Context())

	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	//To change the max size that a form can return (by default its 10MB)
	// r.Body = http.MaxButesReader(w, r.Body, size)

	var form snippetCreateForm

	// An easier way to do the following 4 lines is just to use FormValue or PostFormValue
	// err := r.ParseForm()
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//The get just retuns a string, we would have to convert non-string types from strings.
	// We could have used r.Form instead.
	// r.PostForm handles POST, PATCH, and PUT, whereas the r.Form accepts all types of requests, and any request body, and query string parameters.//

	//Get returns just the first value of a form field, so checkboxes with multiple options are an issue//
	// To make it work, we use something like the following//

	// for i, item := range r.PostForm["items"] {
	// fmt.Fprintf(w, "%d: Item %s\n", i, item)
	// }
	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// form := snippetCreateForm{
	// 	Title:   r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expires,
	// }

	// err = app.formDecoder.Decode(&form, r.PostForm)
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// }

	// Because the Validator type is embedded by the snippetCreateForm struct,
	// we can call CheckField() directly on it to execute our validation checks.
	// CheckField() will add the provided key and error message to the
	// FieldErrors map if the check does not evaluate to true.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
