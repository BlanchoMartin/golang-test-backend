// main.go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "backend/models"

    "github.com/gorilla/mux"
)

var db *gorm.DB

// Initialize database
func initDB() {
    var err error
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Migrate the schema for the Brainee model
    db.AutoMigrate(&models.Brainee{})
}

func main() {
    initDB()

    r := mux.NewRouter()

    // Routes for brainees
    r.HandleFunc("/brainees", handleCreateBrainee).Methods("POST")
    r.HandleFunc("/brainees/{braineeId}", handleGetBrainee).Methods("GET")

    log.Println("Server running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

// POST /brainees - Create a new brainee
func handleCreateBrainee(w http.ResponseWriter, r *http.Request) {
    var brainee models.Brainee
    if err := json.NewDecoder(r.Body).Decode(&brainee); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    db.Create(&brainee)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(brainee)
}

// GET /brainees/{braineeId} - Get a specific brainee by ID
func handleGetBrainee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    braineeId, err := strconv.Atoi(vars["braineeId"])
    if err != nil {
        http.Error(w, "Invalid brainee ID", http.StatusBadRequest)
        return
    }

    var brainee models.Brainee
    if err := db.First(&brainee, braineeId).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            http.Error(w, "Brainee not found", http.StatusNotFound)
        } else {
            http.Error(w, "Error retrieving brainee", http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(brainee)
}
