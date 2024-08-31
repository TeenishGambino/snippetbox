package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	"snippetbox.abiral.net/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	//The underscore is just an alias because the mysql folder can't be seen by the compiler
	// the driver's init function will register itslef with the database/sql package
	// It will happen in run time
	// So we use the alias.
	_ "github.com/go-sql-driver/mysql"
)

// This only is good because we have the handlers in the same package
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	//Define a new command line flag for the MYSQL DSN String
	//It is the password of the database, may not be the best to have it out in the open like this//
	dsn := flag.String("dsn", "web:Hallo123@@/snippetbox?parseTime=true", "MySQL data source name")

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

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	//Dependencies
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// We create a new server so that we can customize it
	//We want to make use of our errorLogger, otherwise the listenAndServe using the default
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// You have to dereference the value because the flag parser just has the location of it and not the value itself.
	// So does that mean parse just keeps in the memory, in a temporary file? It stores it directly in the memory//
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
