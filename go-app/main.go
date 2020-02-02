package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
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

var db *sql.DB
var num int

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/gate", PostHandler)

	// search db if authenticated.
	// if authenticated, send request to gate.

	cfg := mysql.Cfg("qrentry-makeharvard-2020:us-east4:qrentry", "root", "makeharvard2020")
	cfg.DBName = "qrEntryDatabase"
	db, err := mysql.DialCfg(cfg)
	if err != nil {
		//
		fmt.Println(err)
	}

	rows, err := db.Query(`SELECT * FROM GATE`)
	if err != nil {
		//
		fmt.Println(err)
	}
	rows.Scan(&num)
	fmt.Println(num)
	fmt.Println(db)
	fmt.Println(db.Ping())
	fmt.Println(db.Query(`SELECT * FROM AUTHENTICATED_USERS`))

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
