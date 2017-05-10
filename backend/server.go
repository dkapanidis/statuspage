package main

import (
    "fmt"
    "os"
    "net/http"
    "log"
    "math/rand"
    "time"
    "encoding/json"
)

var coin int
var target_url string = os.Getenv("TARGET_URL")
var d Data

type Payload struct {
    Stuff Data
}
type Data struct {
    StatusCodes StatusCodes
}
type StatusCodes map[string]int
type Vegetables map[string]int

var coins Data

func serveHTTP(w http.ResponseWriter, r *http.Request) {
  if (coin == 1) {
    w.Write([]byte("200 - You are a lucky lucky guy!"))
  } else {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 - Bad luck, try flipping the coin again!"))
  }
}

func metrics(w http.ResponseWriter, r *http.Request) {
  response, err := getJsonResponse();
  if err != nil {
      panic(err)
  }
  fmt.Fprintf(w, string(response))
}

func getJsonResponse()([]byte, error) {
  p := Payload{d}
    return json.MarshalIndent(p, "", "  ")
}

var netClient = &http.Client{
  Timeout: time.Second * 10,
}

func doSomething(s string) {
  fmt.Println("doing something", s)
}

func polling() {
  for {
    <-time.After(1 * time.Second)
    stressTest()
  }
}

func stressTest() {
  statusCodes := make(map[string]int)
  statusCodes["20x"] = 0
  statusCodes["50x"] = 0
  for i := 1; i <= 100; i++ {
    response, _ := netClient.Get(target_url)
    if (response.StatusCode == 500) {
      statusCodes["20x"]++
    } else {
      statusCodes["50x"]++
    }

  }
  jsn, _ := json.MarshalIndent(statusCodes, "", "  ")
  fmt.Println(string(jsn))
}

func main() {
    http.HandleFunc("/", serveHTTP)
    http.HandleFunc("/metrics", metrics)
    rand.Seed( time.Now().UnixNano())
    coin = rand.Intn(2)
    statusCodes := make(map[string]int)
    statusCodes["20x"] = 25
    statusCodes["50x"] = 10
    d = Data{statusCodes}

    fmt.Println(coin)
    fmt.Println("TARGET_URL: ", target_url)
    response, foo := netClient.Get(target_url)
    fmt.Println(response.StatusCode)
    fmt.Println(response)
    fmt.Println(foo)
    go polling()
    err := http.ListenAndServe(":80", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
