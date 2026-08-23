# Akritas H6 Demo Fixture

This intentionally defective Go service is used by the H6 backend demo.

Reproduce the error:

```sh
go run .
curl "http://localhost:18080/checkout?customer_id=abc"
curl "http://localhost:18080/checkout?customer_id=abc"
curl "http://localhost:18080/checkout?customer_id=abc"
```

The service panics because `discountCode` slices `customerID[:8]` without
checking length. The regression test in `main_test.go` documents the expected
fix and starts red.
