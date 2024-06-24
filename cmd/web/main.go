package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// This only is good because we have the handlers in the same package
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

//if  you have multiple packages:
// have a config package that has an application structure like above
// then you would something similar to below, but instead you would pass the variable inside the paramter
//Eg: app = &config.Application{...}
//mux.Handle("/", examplePackage.ExampleHandler(app))

func main() {

	//This is just reading from the terminal, you have to add -addr and value ':portNumber' for it
	//It's like scanf in C.
	// Seccond paramter is the default value
	// The third parameter just describes what the flag is for.
	addr := flag.String("addr", ":4000", "HTTP Network address")
	//Flag has Into and Bool, Float64, etc that work similarly, excpet they convert to appropriate types//
	// Doing go run ./cmd/web -help will return the third parameter and the default value//

	//This does the parsing, and sets the value to addr.
	// You need to call this before using the variable//
	flag.Parse()

	// For preexisting variables we could do something like this//
	// flag.StringVar(&addr, "addr", ":4000", "HTTP network address")

	// This creates a logger for writing messages that relate to information.
	//The second parameter is just the prefix
	// | is a bitwise
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Use log.Llongfile instead of Lshortfile for the full file path//
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

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
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// We create a new server so that we can customize it
	//We want to make use of our errorLogger, otherwise the listenAndServe using the default
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// You have to dereference the value because the flag parser just has the location of it and not the value itself.
	// So does that mean parse just keeps in the memory, in a temporary file? It stores it directly in the memory//
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
