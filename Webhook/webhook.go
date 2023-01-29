package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
)

// port - default port to start application on
const port = ":8443"

type WebhookRequest struct {
    Name        string
    Destination string
}

func main() {

    dispatcher := &Dispatcher{
        client:       &http.Client{},
        destinations: make(map[string]string),
        mu:           &sync.Mutex{},
    }

    // preparing HTTP server
    srv := &http.Server{Addr: port, Handler: http.DefaultServeMux}

    // webhook registration handler
    http.HandleFunc("/webhooks", func(resp http.ResponseWriter, req *http.Request) {
        dec := json.NewDecoder(req.Body)
        var wr WebhookRequest
        err := dec.Decode(&wr)
        if err != nil {
            resp.WriteHeader(http.StatusBadRequest)
            return
        }

        dispatcher.add(wr.Name, wr.Destination)
    })

    // start dispatching webhooks
    go dispatcher.Start()

    fmt.Printf("Create webhooks on http://localhost%s/webhooks \n", port)
    // starting server
    err := srv.ListenAndServe()

    if err != http.ErrServerClosed {
        log.Fatalf("listen: %s\n", err)
    }
}

type Dispatcher struct {
    client       *http.Client
    destinations map[string]string
    mu           *sync.Mutex
}

func (d *Dispatcher) Start() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            d.dispatch()
        }
    }
}

func (d *Dispatcher) add(name, destination string) {
    d.mu.Lock()
    d.destinations[name] = destination
    d.mu.Unlock()
}

func (d *Dispatcher) dispatch() {
    d.mu.Lock()
    defer d.mu.Unlock()
    for user, destination := range d.destinations {
        go func(user, destination string) {
            req, err := http.NewRequest("POST", destination, bytes.NewBufferString(fmt.Sprintf("Hello %s, current time is %s", user, time.Now().String())))
            if err != nil {
                // probably don't allow creating invalid destinations
                return
            }

            resp, err := d.client.Do(req)
            if err != nil {
                // should probably check response status code and retry if it's timeout or 500
                return
            }

            fmt.Printf("Webhook to '%s' dispatched, response code: %d \n", destination, resp.StatusCode)

        }(user, destination)
    }
}