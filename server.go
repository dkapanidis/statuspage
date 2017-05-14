package main

import (
    "net/http"
    "log"
    "time"
    "encoding/json"
    "os"
)

var target_url string = os.Getenv("TARGET_URL")

/** Agggregated snapshot information of one second **/
type Snapshot struct {
  Value float64 `json:"value"`
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
  var counter float64 = 0
  for i := 1; i <= 100; i++ {
    response, _ := netClient.Get(target_url)
    if (response != nil && response.StatusCode == 200) {
      counter++
    }
  }
  snapshot := Snapshot{Value:counter/100, Date: time.Now()}
  snapshots = append(snapshots, snapshot)
}

func main() {
  // Handlers
  fs := http.FileServer(http.Dir("static"))
  http.Handle("/", fs)
  http.HandleFunc("/metrics", metrics)

  // Activate Polling
  go polling()

  // Start Listener
  err := http.ListenAndServe(":9000", nil) // set listen port
  if err != nil {
      log.Fatal("ListenAndServe: ", err)
  }
}
