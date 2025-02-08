create-migrations:
	migrate create -ext sql -dir ./internal/ifrastructure/db/migrations -seq add_refresh_tab

migrateup:
	migrate -path ./internal/app/repository/migrations -database 'postgres://postgres:postgres@localhost:5432/auth-service?sslmode=disable' up

migratedown:
	migrate -path ./internal/app/repository/migrations -database 'postgres://postgres:postgres@localhost:5432/auth-service?sslmode=disable' down

test-mock:
	mockgen -source=internal/app/service/service.go -destination=internal/app/service/mocks/mock_service.go -package=mocks

gen-docs:
	swag init -g ./cmd/main.go -o ./docs