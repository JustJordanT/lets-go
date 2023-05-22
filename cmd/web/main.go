package main

import (
	// New import "flag"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.01 (Project Setup and Creating // a Module) so that the import statement looks like this:
	// "{your-module-path}/internal/models". If you can't remember what module path you // used, you can find it at the top of the go.mod file.
	"github.com/justjordant/lets-go/internal/models"

	_ "github.com/go-sql-driver/mysql" // New import
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "pqlzpezuck24wv6bxvi4:pscale_pw_tcD6Qs4gTyWenFCsvmiQ3VJIFf2Iey1ShWkVQxIdh1z@tcp(aws.connect.psdb.cloud)/snippet-cloud?tls=true&parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy I've put the code for creating a connection // pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err1 := openDB(*dsn)
	if err1 != nil {
		errorLog.Fatal(err1)
	}
	// We also defer a call to db.Close(), so that the connection pool is closed // before the main() function exits.
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{
			DB: db,
		},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		Handler: app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
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