package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Address representation.
type Address struct {
	IP net.IP `json:"ip"`
}

// parseIP returns a valid IP address from an x-forwarded-for header.
func parseIP(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info(r.Header.Get("X-Forwarded-For"))
		ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
		if ip == nil {
			http.Error(w, "Error parsing IP address", http.StatusInternalServerError)
			return
		}

		// Plaintext by default, JSON with query param
		if _, ok := r.URL.Query()["json"]; ok {
			w.Header().Set("Content-Type", "application/json")

			// ip will always be encodeable
			_ = json.NewEncoder(w).Encode(Address{IP: ip})

			return
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, ip)
	}
}
