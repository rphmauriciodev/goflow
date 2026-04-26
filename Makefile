include .env
export

# Roda a aplicação com os logs identados
run:
	go run cmd/goflow/main.go | jq

# Roda as migrações
migrate-up:
	migrate -path db/migrations -database "$(DATABASE_URL)?sslmode=disable" up

# Reverte as migrações
migrate-down:
	migrate -path db/migrations -database "$(DATABASE_URL)?sslmode=disable" down

# Exibe o resumo do sistema
dashboard:
	go run cmd/dashboard/main.go

# Bônus: Watch mode (roda o dashboard a cada 2 segundos)
watch-dashboard:
	watch -n 2 go run cmd/dashboard/main.go