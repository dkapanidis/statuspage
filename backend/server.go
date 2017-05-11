package main

import (
    "fmt"
    "net/http"
    "log"
    "time"
    "encoding/json"
)

var coin int
// var target_url string = os.Getenv("TARGET_URL")
var target_url string = "http://192.168.0.155:8080"
var d Data

type Payload struct {
    Stuff Data
}

var output = make([]map[string]int, 100)

/* */
type Data struct {
    StatusCodes StatusCodes
}

type StatusCodes map[string]int
type Vegetables map[string]int

type Snapshots []Snapshot

type Snapshot struct {
  R200 int `json:"r200"`
  R500 int `json:"r500"`
}

var coins Data
var snapshots Snapshots

func metrics(w http.ResponseWriter, r *http.Request) {
  // response, err := getJsonResponse();
  // if err != nil {
      // panic(err)
  // }
  json.NewEncoder(w).Encode(snapshots)

  // w.Header().Set("Content-Type", "application/json")
  // v,err := json.Marshal(fmt.Sprintf("%+v", timeseries))
  // fmt.Printf("TS %+v\n", v)
  // fmt.Printf("ERR %+v\n", err)
  // // fmt.Fprintf(w, fmt.Sprintf("%+v\n", timeseries))
  // fmt.Fprintf(w, string(v))
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
  snapshot := Snapshot{R200:0, R500:0}
  for i := 1; i <= 100; i++ {
    response, _ := netClient.Get(target_url)
    if (response != nil && response.StatusCode == 500) {
      snapshot.R200++
    } else {
      snapshot.R500++
    }
  }
  snapshots = append(snapshots, snapshot)
  fmt.Printf("%+v\n", snapshot)
  fmt.Printf("%+v\n", snapshots)

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
