package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	listen := flag.String("listen", ":9911", "listen `address`, like \":80\" or \"localhost:9999\"")
	target := flag.String("alerts.url", "http://admin:admin@emergency:23233/alerts/list", "`URL` of /alerts/list endpoint")
	metric := flag.String("heartbeat.metric", "obs_alert_heartbeat_time_seconds", "metric `name` for emergency alert heartbeat timestamp")
	flag.Parse()

	client := &http.Client{}
	srv := &http.Server{
		Addr: *listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error
			var code int
			defer func() {
				if code == 0 {
					code = http.StatusOK
				} else {
					msg := http.StatusText(code)
					if err != nil {
						msg = err.Error()
					}
					http.Error(w, msg, code)
				}
				log.Printf("%q %d %q %q", r.RemoteAddr, code, r.Method, r.URL.Path)
			}()
			if r.URL.Path != "/metrics" {
				code = http.StatusNotFound
				return
			}
			req, err := http.NewRequest("POST", *target, nil)
			if err != nil {
				code = http.StatusInternalServerError
				return
			}
			req = req.WithContext(r.Context())
			resp, err := client.Do(req)
			if err != nil {
				code = http.StatusBadGateway
				return
			}
			var alerts struct {
				LastHeartbeat float64 `json:"last_heartbeat"`
			}
			err = json.NewDecoder(resp.Body).Decode(&alerts)
			if err != nil {
				code = http.StatusBadGateway
				return
			}
			fmt.Fprintf(w, "%s %f\n", *metric, alerts.LastHeartbeat)
		}),
	}
	log.Fatal(srv.ListenAndServe())
}
