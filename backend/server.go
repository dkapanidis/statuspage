package main

import (
    "net/http"
    "log"
    "time"
    "encoding/json"
)

// var target_url string = os.Getenv("TARGET_URL")
var target_url string = "http://192.168.0.155:8080"

/** Agggregated snapshot information of one second **/
type Snapshot struct {
  R200 int `json:"r200"`
  R500 int `json:"r500"`
  Date time.Time `json:"date"`
}

type Snapshots []Snapshot


var snapshots Snapshots

func metrics(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(snapshots)
}

var netClient = &http.Client{
  Timeout: time.Second * 10,
}

// Each second run stressTest
func polling() {
  for {
    <-time.After(1 * time.Second)
    stressTest()
  }
}

// Run 100 cycles and collect snapshot
func stressTest() {
  snapshot := Snapshot{R200:0, R500:0, Date: time.Now()}
  for i := 1; i <= 100; i++ {
    response, _ := netClient.Get(target_url)
    if (response != nil && response.StatusCode == 500) {
      snapshot.R200++
    } else {
      snapshot.R500++
    }
  }
  snapshots = append(snapshots, snapshot)
}

func main() {
  // Handlers
  http.HandleFunc("/", metrics)

  // Activate Polling
  go polling()

  // Start Listener
  err := http.ListenAndServe(":8081", nil) // set listen port
  if err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}
