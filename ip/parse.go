package ip

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

// Parse a valid IP address from an x-forwarded-for header.
func Parse(w http.ResponseWriter, r *http.Request) {
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip == nil {
		http.Error(w, "Error parsing IP address", http.StatusInternalServerError)
		return
	}

	// Plaintext by default, JSON with query param
	if _, ok := r.URL.Query()["json"]; ok {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(Address{IP: ip})

		if err != nil {
			http.Error(w, "Invalid Request", http.StatusInternalServerError)
			return
		}

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, ip)
}
