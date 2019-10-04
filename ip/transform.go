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

// Transform a header into a IP address response
func Transform(w http.ResponseWriter, r *http.Request) {
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip == nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Plaintext by default, JSON with query param
	if _, ok := r.URL.Query()["json"]; ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Address{ip.String()})
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, ip.String())
}
