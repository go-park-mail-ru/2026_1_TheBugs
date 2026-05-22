package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"
)

var AllowYookassaIPs = []string{
	"185.71.76.0/27",
	"185.71.77.0/27",
	"77.75.153.0/25",
	"77.75.156.11",
	"77.75.156.35",
	"77.75.154.128/25",
	"2a02:5180::/32",
}

func IPFilterMiddleware(next http.Handler, allowedCIDRs []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIP := r.RemoteAddr
		log.Printf("Initial remote IP: %s", remoteIP)
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			log.Printf("Using X-Real-IP header: %s", realIP)
			remoteIP = realIP
		}

		var host string
		if strings.Contains(remoteIP, ":") {
			var err error
			host, _, err = net.SplitHostPort(remoteIP)
			if err != nil {
				http.Error(w, "Invalid remote IP address", http.StatusBadRequest)
				return
			}
		} else {
			host = remoteIP
		}

		if !IsIPAllowed(host, allowedCIDRs) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IsIPAllowed(ip string, allowedCIDRs []string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, cidr := range allowedCIDRs {
		_, allowedNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if allowedNet.Contains(parsedIP) {
			return true
		}
	}
	return false
}
