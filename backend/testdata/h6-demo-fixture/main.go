package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/checkout", checkoutHandler)
	addr := ":18080"
	log.Printf("h6 demo fixture listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	_, _ = fmt.Fprintln(w, discountCode(customerID))
}

func discountCode(customerID string) string {
	prefix := os.Getenv("AKRITAS_DEMO_DISCOUNT_PREFIX")
	if prefix == "" {
		prefix = "DEMO"
	}
	return prefix + "-" + customerID[:8]
}
