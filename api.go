package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "sync"
)

var (
    lock sync.RWMutex
    monitoringDataCache map[string]interface{}
)

func init() {
    data := os.Getenv("MONITORING_DATA")
    if data != "" {
        json.Unmarshal([]byte(data), &monitoringDataCache)
    }
    if monitoringDataCache == nil {
        monitoringDataCache = make(map[string]interface{})
    }
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"message": "Everything is OK!"}
    json.NewEncoder(w).Encode(response)
}

func RetrieveMonitoringData(w http.ResponseWriter, r *http.Request) {
    lock.RLock()
    defer lock.RUnlock()
  
    if monitoringDataCache == nil {
        http.Error(w, "No monitoring data available", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(monitoringDataCache)
}

func UpdateMonitoringData(w http.ResponseWriter, r *http.Request) {
    var newData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&newData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    lock.Lock()
    monitoringDataCache = newData
    lock.Unlock()
    
    updatedDataJSON, err := json.Marshal(newData)
    if err != nil {
        http.Error(w, "Failed to update monitoring data", http.StatusInternalServerError)
        return
    }
    os.Setenv("MONITORING_DATA", string(updatedDataJSON))
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(newData)
}

func main() {
    http.HandleFunc("/health", HandleHealthCheck)
    http.HandleFunc("/data", RetrieveMonitoringData)
    http.HandleFunc("/data/update", UpdateMonitoringData)

    port := os.Getenv("PORT")
    if port == "" {
        log.Fatal("$PORT must be set")
    }

    log.Printf("Starting server on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}