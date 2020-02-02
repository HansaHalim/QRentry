package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func openHandler(w http.ResponseWriter, r *http.Request) {
	if accessGranted {
		//Give get response code
		w.WriteHeader(http.StatusOK)
	} else {
		// bad response
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//PostHandler handles gate requests
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		item := &userRequest{}

		item.email = r.Header.Get("email")
		item.gateID = r.Header.Get("gateID")

		fmt.Fprint(w, fmt.Sprintf("email: %s, gateID: %s, dbRow: %d\n", item.email, item.gateID, num))
		query := fmt.Sprintf("SELECT * FROM qrEntryDatabase.AUTHENTICATED_USERS WHERE gateID=%s AND email='%s'", item.gateID, item.email)
		searchRow := db.QueryRow(query)

		var count int       // This is to take in gateID
		var anything string // This takes the email
		searchRow.Scan(&count, &anything)
		fmt.Fprint(w, "\n")
		fmt.Fprint(w, fmt.Sprintf("GateID: %d, Email: %s\n", count, anything))
		if strconv.Itoa(count) == item.gateID {
			fmt.Fprint(w, "Access Granted!\n")
			accessGranted = true
			delay, _ := time.ParseDuration("5s")
			time.Sleep(delay)
			accessGranted = false
			// Call Function to activate gate here
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func DB() *sql.DB {
	var (
		connectionName = "qrentry-makeharvard-2020:us-east4:qrentry"
		user           = "root"
		dbName         = "qrEntryDatabase"
		password       = "makeharvard2020"
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
var accessGranted bool

func main() {
	accessGranted = false
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/gate", PostHandler) // Search db if user is authenticated
	http.HandleFunc("/open", openHandler) // allow gate to be opened if user is authenticated

	db = DB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
