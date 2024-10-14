package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "sync"
)

var (
    dataLock            sync.RWMutex
    apiMonitoringData   map[string]interface{}
)

func init() {
    initialData := os.Getenv("MONITORING_DATA")
    if initialData != "" {
        json.Unmarshal([]byte(initialData), &apiMonitoringData)
    }
    if apiMonitoringData == nil {
        apiMonitoringData = make(map[string]interface{})
    }
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
    statusResponse := map[string]string{"message": "Everything is OK!"}
    json.NewEncoder(w).Encode(statusResponse)
}

func getMonitoringData(w http.ResponseWriter, r *http.Request) {
    dataLock.RLock()
    defer dataLock.RUnlock()
  
    if apiMonitoringData == nil {
        http.Error(w, "No monitoring data available", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(apiMonitoringData)
}

func updateMonitoringData(w http.ResponseWriter, r *http.Request) {
    var newMonitoringData map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&newMonitoringData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    dataLock.Lock()
    apiMonitoringData = newMonitoringData
    dataLock.Unlock()
    
    updatedDataJSON, err := json.Marshal(newMonitoringData)
    if err != nil {
        http.Error(w, "Failed to update monitoring data", http.StatusInternalServerError)
        return
    }
    os.Setenv("MONITORING_DATA", string(updatedDataJSON))
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(newMonitoringData)
}

func main() {
    http.HandleFunc("/health", handleHealthCheck)
    http.HandleFunc("/data", getMonitoringData)
    http.HandleFunc("/data/update", updateMonitoringData)

    serverPort := os.Getenv("PORT")
    if serverPort == "" {
        log.Fatal("$PORT must be set")
    }

    log.Printf("Starting server on port %s", serverPort)
    log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}