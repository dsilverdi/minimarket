#!/bin/bash
cd user/api/
protoc --go_out=. --go-grpc_out=. user.proto 
cd .. && cd ..
cd product/api/
protoc --go_out=. --go-grpc_out=. product.proto 