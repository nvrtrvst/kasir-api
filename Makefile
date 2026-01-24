# Perintah untuk update swagger lalu jalankan server
run:
	@C:\Users\End\go\bin\swag init -g cmd/server/main.go --parseInternal --dir ./
	@go run cmd/server/main.go
