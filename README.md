# Start the port service

## On client service
```golang
go run cmd/main.go
```

## On server
```golang
go run cmd/main.go
```

## Run the request form Postman or any other service

GET http://localhost:3000/ports  --to list all ports
GET http://localhost:3000/ports/{id} --to list the port
POST http://localhost:3000/ports --to upsert ports

# Docker
## Run in server folder

docker-compose up --build client
## Run the requests as showed earlier