# Wasmtime-go practice

## Quick start

```
go get
make start-docker-test
make build-wasm
go run main.go
docker exec test psql -U test -d test -c "SELECT * FROM person"
```