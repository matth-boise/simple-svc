package main

import (
  "os"
  "io/ioutil"
  "fmt"
  "log"
  "strings"
  "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  service, exists := os.LookupEnv("SERVICE_NAME")
  if !exists {
    service = "unknown-service"
  }
  ip, exists := os.LookupEnv("IP_ADDRESS")
  if !exists {
    ip = "unknown-IP"
  }
  log.Printf("%s called with Host=%s", service, r.Host)
  
  // printing to the http.ResponseWriter returns text to http caller
  // TODO: add env POD_NAME

  if s := strings.Split(r.URL.Path, "/"); len(s) > 2 {
    nextService := s[1]
    nextPathSlice := s[2:]
    nextPath := strings.Join(nextPathSlice, "/")
    url := "http://" + nextService + "/" + nextPath
    //fmt.Fprintf(w, "%s: TODO: call http://%s/%s and return response!\n", service, nextService, nextPath)
    response, err := http.Get(url)
    if err != nil {
      log.Printf("ERR: GET(%s)\n", url)
      log.Fatal(err)
    } else {
      defer response.Body.Close()
      responseBody, _ := ioutil.ReadAll(response.Body)
      responseText := string(responseBody)
      fmt.Fprintf(w, "%s-%s: GET %s returned code=%d body=%s", service, ip, url, response.StatusCode, responseText)
      log.Printf("%s-%s: call to GET %s returned code=%d body=%s", service, ip, url, response.StatusCode, responseText)
    }
  }  else {
    response := service + "-" + ip + ": path=" + r.URL.Path
    fmt.Fprintf(w, "%s\n", response)
    log.Printf("%s-%s: return \"%s\"\n", service, ip, response)
  }
}

func main() {
  port, portSpecified := os.LookupEnv("SERVICE_PORT")
  if !portSpecified {
    port = "8000"
  }

  // specify path
  http.HandleFunc("/", handler)
  // listen on port
  log.Fatal(http.ListenAndServe(":" + port, nil))
}
