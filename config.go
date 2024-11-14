package main

import (
    "fmt"
    "os"
    "sync"

    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
        return
    }

    apiKey := os.Getenv("API_KEY")
    dbSettings := loadDatabaseSettings()

    fmt.Printf("API Key: %s\n", dbSettings.apiKey)

    batchProcessAPICalls(apiKey)
}

func batchProcessAPICalls(apiKey string) {
    items := []string{"item1", "item2", "item3"}

    var wg sync.WaitGroup

    batchSize := 10 
    for i := 0; i < len(items); i += batchSize {
        end := i + batchSize
        if end > len(items) {
            end = len(items)
        }

        wg.Add(1)
        go func(batch []string) {
            defer wg.Done()
            for _, item := range batch {
                fetchData(apiKey, item)
            }
        }(items[i:end])
    }

    wg.Wait() 
}

func fetchData(apiKey, item string) {
    fmt.Printf("Fetching data for %s with API Key %s\n", item, apiKey)
}

func loadDatabaseSettings() *dbConfig {
    return &dbConfig{
        apiKey:      os.Getenv("API_KEY"),
        dbHost:      os.Getenv("DB_HOST"),
        dbPort:      os.Getenv("DB_PORT"),
        dbUser:      os.Getenv("DB_USER"),
        dbPassword:  os.Getenv("DB_PASSWORD"),
        dbName:      os.Getenv("DB_NAME"),
        logLevel:    os.Getenv("LOG_LEVEL"),
    }
}

type dbConfig struct {
    apiKey      string
    dbHost      string
    dbPort      string
    dbUser      string
    dbPassword  string
    dbName      string
    logLevel    string
}