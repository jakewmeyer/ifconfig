package ip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Address representation
type Address struct {
	IP string `json:"ip"`
}

// Get a valid IP address from an x-forwarded-for header
func Get(w http.ResponseWriter, r *http.Request) {
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip == nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Plaintext by default, JSON with query param
	if _, ok := r.URL.Query()["json"]; ok {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Address{ip.String()})
		if err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s", ip.String())
}
