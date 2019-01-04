package main

import (
  "os"
  "io/ioutil"
  "fmt"
  "log"
  "strings"
  "sort"
  "net/http"
)

func getServiceAndIP() (string, string) {
  service, exists := os.LookupEnv("SERVICE_NAME")
  if !exists {
    service = "unknown-service"
  }
  ip, exists := os.LookupEnv("IP_ADDRESS")
  if !exists {
    ip = "unknown-IP"
  }
  return service, ip
}

// if headerNames header exist in srcRequest.Headers, copy them to to destRequest.Headers
func copyHeaders(headerNames []string, srcRequest *http.Request, destRequest *http.Request) {
  for _, headerName := range headerNames {
    // TODO: make sure headerName exists in srcRequest.Header
    if headerValue := srcRequest.Header.Get(headerName); headerValue != "" {
      destRequest.Header.Set(headerName, headerValue)
    }
  }
}


func handler(w http.ResponseWriter, r *http.Request) {
  service, ip := getServiceAndIP()

  log.Printf("%s called with Host: %s  and other headers:", service, r.Host)
  var headers []string
  for h := range r.Header {
    headers = append(headers, h)
  }
  sort.Strings(headers)

  for _, h := range headers {
     log.Printf("      %s: %s", h, r.Header[h])
   }

  // printing to the http.ResponseWriter returns text to http caller
  // TODO: add env POD_NAME

  if s := strings.Split(r.URL.Path, "/"); len(s) > 2 {
    nextService := s[1]
    nextPathSlice := s[2:]
    nextPath := strings.Join(nextPathSlice, "/")
    url := "http://" + nextService + "/" + nextPath
    //fmt.Fprintf(w, "%s: TODO: call http://%s/%s and return response!\n", service, nextService, nextPath)

    client := &http.Client{
      CheckRedirect: nil,
    }
    //response, err := http.Get(url)
    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
      log.Printf("ERR: http.NewRequest(\"GET\", url=%s)\n", url)
      log.Fatal(err)
    } else {
      headerNames := []string{"X-Request-Id", "X-B3-Traceid", "X-B3-Spanid", "X-B3-Sampled", "X-B3-Parentspanid", "X-B3-Flags"}
      copyHeaders(headerNames, r, request)
      response, err := client.Do(request)
      //defer response.Body.Close()
      if err != nil {
        log.Printf("ERR: GET(%s)\n", url)
        log.Fatal(err)
      } else {
        responseBody, _ := ioutil.ReadAll(response.Body)
        responseText := string(responseBody)
        fmt.Fprintf(w, "%s-%s: GET %s returned code=%d body=%s", service, ip, url, response.StatusCode, responseText)
        log.Printf("%s-%s: call to GET %s returned code=%d body=%s", service, ip, url, response.StatusCode, responseText)
      }
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
