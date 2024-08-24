package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

func serveHomePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: HomePage")
}

func initializeServer() {
    portEnv := os.Getenv("PORT")
    if portEnv == "" {
        portEnv = "8080"
    }

    serverAddress := ":" + portEnv
    http.HandleFunc("/", serveHomePage)
    log.Printf("Starting server on port %s\n", portEnv)
    log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    initializeServer()
}