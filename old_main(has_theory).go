package old_main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
// r is the request that the user could be sending, i.e Post or something//
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))

}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	// w.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {

	//Using Http.MethodPost and http.method<> is better than using string because
	// It prevents runtime errors due to typos.
	if r.Method != http.MethodPost {

		// Use the Header().Set() method to add an 'Allow: POST' header to the
		// response header map. The first parameter is the header name, and
		// the second parameter is the header value.
		w.Header().Set("Allow", http.MethodPost)

		// Set a new cache-control header. If an existing "Cache-Control" header exists
		// it will be overwritten.
		w.Header().Set("Cache-Control", "public, max-age=31536000")

		// In contrast, the Add() method appends a new "Cache-Control" header and can
		// be called multiple times.
		w.Header().Add("Cache-Control", "public")
		w.Header().Add("Cache-Control", "max-age=31536000")

		// Delete all values for the "Cache-Control" header.
		w.Header().Del("Cache-Control")

		// Retrieve the first value for the "Cache-Control" header.
		w.Header().Get("Cache-Control")

		// Retrieve a slice of all values for the "Cache-Control" header.
		w.Header().Values("Cache-Control")

		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		//This is equivalent to
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		return
	}
	w.Write(([]byte("Create a new snippet...")))
}

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	//Mapping / -> home
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.

	//If we had a lot of urls then we would just have more mux.HandleFunc the below one would
	// just be the same.
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
