package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

//go:embed index.html
var f embed.FS

var DB *sql.DB

type ToDo struct {
	Content string
}

func connectDB() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?multiStatements=true", dbUser, dbPass, dbHost, dbName))

	if err != nil {
		return nil, fmt.Errorf("Error opening database: %v", err)
	}

	return db, nil
}

func migration() error {
	log.Println("Starting migration")
	driver, err := mysql.WithInstance(DB, &mysql.Config{})

	if err != nil {
		return fmt.Errorf("Error creating driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)

	if err != nil {
		return fmt.Errorf("Error creating migration instance: %v", err)
	}

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Error running migration: %v", err)
	}

	log.Println("Migration completed")
	return nil
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT content FROM todo")
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer rows.Close()

	todos := []ToDo{}
	for rows.Next() {
		var t ToDo
		err := rows.Scan(&t.Content)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		todos = append(todos, t)
	}

	tmpl, _ := template.ParseFS(f, "index.html")
	err = tmpl.Execute(w, todos)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	rows, err := DB.Query(fmt.Sprintf("INSERT INTO todo (content) VALUES ('%s')", content))
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer rows.Close()
	http.Redirect(w, r, "/", 302)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [command]")
		os.Exit(0)
	}

	log.Println("Start")

	db, err := connectDB()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	DB = db
	defer DB.Close()

	switch os.Args[1] {
	case "start":
		http.HandleFunc("/", RootHandler)
		http.HandleFunc("/todo", PostHandler)
		http.HandleFunc("/health_check", HealthCheckHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	case "migrate":
		err := migration()
		if err != nil {
			log.Fatalf("Error running migration: %v", err)
		}
	}
}
