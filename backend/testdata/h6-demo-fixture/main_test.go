package main

import "testing"

func TestDiscountCodeHandlesShortCustomerID(t *testing.T) {
	t.Setenv("AKRITAS_DEMO_DISCOUNT_PREFIX", "H6")
	got := discountCode("abc")
	if got != "H6-abc" {
		t.Fatalf("discountCode() = %q, want %q", got, "H6-abc")
	}
}
