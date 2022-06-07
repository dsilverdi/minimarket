# Mini Market

set up the whole server module for mini marketplace platform

## Project Structure
```
.
├── cmd
|   ├── gateway
|   |   └── main.go     // Main App for Gateway Service
|   ├── product
|   |   └── main.go     // Main App for Product Service
|   └── user
|       └── main.go     // Main App for User Service
├── gateway             // Packages used inside the project, including packages for crud, service (facade) and business logic
│   ├── api
│   ├── client       
│   └── mock   
├── product             // Packages used inside the project, including packages for grpc, service (facade) and business logic
│   ├── api
│   ├── database       
│   └── mock
├── gateway             // Packages used inside the project, including packages for grpc, service (facade) and business logic
│   ├── api
│   ├── database       
│   └── mock              
├── pkg                 // Additional Package            
├── docker-compose.yml
├── Dockerfile
├── generate-pb.sh
├── uuid.go
├── go.mod
├── go.sum
└── README.md
```


## Technical Constrainst
* Go 1.14
* MySQL 5.7
* Redis 4.0+

## Depedencies
* Docker
* Docker Compose

## How to run
to run this service you should have Docker installed beforehand. 

build service
```sh-session
$ docker compose build
```

to run this service
```sh-session
$ docker compose up
```