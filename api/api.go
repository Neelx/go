package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Item represents a data item in our system
type Item struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// App encapsulates the application
type App struct {
	DB *sql.DB
}

// Initialize sets up the database connection and creates tables if needed
func (a *App) Initialize(dbPath string) error {
	var err error
	a.DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Create table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		value TEXT NOT NULL
	);`

	_, err = a.DB.Exec(query)
	return err
}

// GetItems returns all items from the database
func (a *App) GetItems(w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.Query("SELECT id, name, value FROM items")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Value); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		items = append(items, i)
	}

	respondWithJSON(w, http.StatusOK, items)
}

// GetItem returns a specific item by ID
func (a *App) GetItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing item ID")
		return
	}

	var i Item
	err := a.DB.QueryRow("SELECT id, name, value FROM items WHERE id = ?", id).Scan(&i.ID, &i.Name, &i.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, i)
}

// CreateItem adds a new item to the database
func (a *App) CreateItem(w http.ResponseWriter, r *http.Request) {
	var i Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	stmt, err := a.DB.Prepare("INSERT INTO items(name, value) VALUES(?, ?)")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(i.Name, i.Value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	i.ID = int(id)
	respondWithJSON(w, http.StatusCreated, i)
}

// UpdateItem updates an existing item
func (a *App) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "", err := strconv.Atoi(id); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var i Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&i); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	stmt, err := a.DB.Prepare("UPDATE items SET name = ?, value = ? WHERE id = ?")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(i.Name, i.Value, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, i)
}

// DeleteItem removes an item from the database
func (a *App) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing item ID")
		return
	}

	stmt, err := a.DB.Prepare("DELETE FROM items WHERE id = ?")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Helper functions for HTTP responses
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error encoding response"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// SetupRoutes configures the API endpoints
func (a *App) SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	
	// API endpoints
	mux.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			a.GetItems(w, r)
		case http.MethodPost:
			a.CreateItem(w, r)
		default:
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/api/item", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			a.GetItem(w, r)
		case http.MethodPut:
			a.UpdateItem(w, r)
		case http.MethodDelete:
			a.DeleteItem(w, r)
		default:
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// Serve static files for frontend
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	return mux
}

// Run starts the API server
func (a *App) Run(addr string) {
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.SetupRoutes()))
}

func main() {
	app := App{}
	err := app.Initialize("./database.db")
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}

	fmt.Println("API server initialized. Database connected.")
	app.Run(":8080")
}