package main

import (
    "fmt"
    "os"
    "net/http"
    "log"
    "math/rand"
    "time"
)

var coin int
var target_url string = os.Getenv("TARGET_URL")

var coin_ok int = 0
var coin_fail int = 0

func serveHTTP(w http.ResponseWriter, r *http.Request) {
  if (coin == 1) {
    w.Write([]byte("200 - You are a lucky lucky guy!"))
  } else {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 - Bad luck, try flipping the coin again!"))
  }
}

func metrics(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("{coin_ok: " + coin_ok + ", coin_fail: 4}"))
}

var netClient = &http.Client{
  Timeout: time.Second * 10,
}

func main() {
    http.HandleFunc("/", serveHTTP)
    http.HandleFunc("/metrics", metrics)
    rand.Seed( time.Now().UnixNano())
    coin = rand.Intn(2)
    fmt.Println(coin)
    fmt.Println("TARGET_URL: ", target_url)
    response, foo := netClient.Get(target_url)
    fmt.Println(response)
    fmt.Println(foo)
    err := http.ListenAndServe(":80", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
