package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
)

type userRequest struct {
	email  string
	gateID string
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		log.Printf("Could not query db: %v", err)
		http.Error(w, "Internal Error", 500)
		return
	}
	defer rows.Close()

	newRows, err := db.Query("SELECT * FROM qrEntryDatabase.GATE")

	buf := bytes.NewBufferString("Databases:\n")
	fmt.Fprintf(buf, "%s\n", newRows)
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			log.Printf("Could not scan result: %v", err)
			http.Error(w, "Internal Error", 500)
			return
		}
		fmt.Fprintf(buf, "- %s\n", dbName)
	}

	w.Write(buf.Bytes())
	fmt.Fprint(w, "Hello, World in Harvard!")
}

//PostHandler handles gate requests
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		item := &userRequest{}

		item.email = r.Header.Get("email")
		item.gateID = r.Header.Get("gateID")

		// fmt.Sprintf("email: %s, gateID: %s\n", item.email, item.gateID)``
		fmt.Fprint(w, fmt.Sprintf("email: %s, gateID: %s, dbRow: %d\n", item.email, item.gateID, num))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func DB() *sql.DB {
	var (
		connectionName = "qrentry-makeharvard-2020:us-east4:qrentry"
		user           = "root"
		dbName         = "qrEntryDatabase" //os.Getenv("CLOUDSQL_DATABASE_NAME") // NOTE: dbName may be empty
		password       = "makeharvard2020" // NOTE: password may be empty
		socket         = "/cloudsql"
	)
	// connection string format: USER:PASSWORD@unix(/cloudsql/PROJECT_ID:REGION_ID:INSTANCE_ID)/[DB_NAME]
	dbURI := fmt.Sprintf("%s:%s@unix(%s/%s)/%s", user, password, socket, connectionName, dbName)
	conn, err := sql.Open("mysql", dbURI)
	if err != nil {
		panic(fmt.Sprintf("DB: %v", err))
	}
	return conn
}

var db *sql.DB
var num int

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/gate", PostHandler)

	db = DB()

	// search db if authenticated.
	// if authenticated, send request to gate.

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil { // CHANGE THIS BEFORE PUSHING TO GOOGLE CLOUD "localhost:" -> ":"
		log.Fatal(err)
	}
}
