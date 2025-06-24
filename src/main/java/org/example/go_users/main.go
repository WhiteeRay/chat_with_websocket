package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Password string `json:"password"`
}

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=intern_websocket sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }
    defer db.Close()

    r := mux.NewRouter()

    r.HandleFunc("/users", CreateUser).Methods("POST")
    r.HandleFunc("/login", LoginUser).Methods("POST")
    r.HandleFunc("/users", GetUsers).Methods("GET")
    r.HandleFunc("/users/{id}", GetUser).Methods("GET")
    r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
    r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

    log.Println("Server running at http://localhost:9090")
    log.Fatal(http.ListenAndServe(":9090", enableCORS(r)))
}

func enableCORS(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            return
        }
        h.ServeHTTP(w, r)
    })
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)

    err := db.QueryRow("INSERT INTO users(name, password) VALUES($1, $2) RETURNING id", user.Name, user.Password).Scan(&user.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, password FROM users")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        err := rows.Scan(&u.ID, &u.Name, &u.Password)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }

    json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var user User
    err = db.QueryRow("SELECT id, name, password FROM users WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Password)
    if err == sql.ErrNoRows {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var user User
    json.NewDecoder(r.Body).Decode(&user)

    _, err = db.Exec("UPDATE users SET name=$1, password=$2 WHERE id=$3", user.Name, user.Password, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    user.ID = id
    json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("DELETE FROM users WHERE id=$1", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
    var creds User
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    var user User
    err = db.QueryRow("SELECT id, name, password FROM users WHERE name=$1", creds.Name).Scan(&user.ID, &user.Name, &user.Password)
    if err == sql.ErrNoRows {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if user.Password != creds.Password {
        http.Error(w, "Incorrect password", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(user)
}
