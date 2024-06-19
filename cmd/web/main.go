package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	// if the user search /dab/static/something, the
	// stripPrefix will take the /dab/ and replace /static/something with
	// ./ui/static/something, and join them so that we have
	// /dab/ui/static/something
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//The handle functions make use of the Handler interface,
	// normally you would have to do something like:
	/*
		type home struct {}
		func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("This is my home page"))
		}

		mux := http.NewServeMux() mux.Handle("/", &home{})
	*/

	// This is way to bloaty to do for many things,
	// So we make use of http.HandlerFunc(home) which takes the contents in the home function
	// And ServeHTTP just calls it.

	// What we do here is just syntatic sugar so that we don't have to keep typing http.HandlerFunction all the time.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
