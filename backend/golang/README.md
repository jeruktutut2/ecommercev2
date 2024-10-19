# Backend Ecommerce

## install echo
go get github.com/labstack/echo/v4

## install postgres
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool

## install redis
go get github.com/redis/go-redis/v9

## install validator
go get github.com/go-playground/validator/v10

## install testify
go get github.com/stretchr/testify

## install uuid
go get github.com/google/uuid


## test
go test -v tests/integration_tests/features/users/login/services/login_service_test.go  
go test -v tests/unit_tests/features/users/login/services/login_service_test.go  
go test -v tests/api_tests/features/users/login/login_test.go  
chmod +x login_curl.sh  
./login_curl.sh  