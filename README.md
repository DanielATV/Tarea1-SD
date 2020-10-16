# Tarea1-SD

protoc --go_out=plugins=grpc:chat chat.proto
protoc --go_out=plugins=grpc:camion camion.proto

go run hello-world.go
go build hello-world.go
./hello-world


export GOROOT=/usr/local/go ; export GOPATH=$HOME/go ; export GOBIN=$GOPATH/bin ; export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN