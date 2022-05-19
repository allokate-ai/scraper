go install github.com/cespare/reflex@latest

$(go env GOPATH)/bin/reflex --start-service -r ".*\.go" -R ".*_test\.go" go -- run cmd/main.go

