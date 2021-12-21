FROM golang
# Copy our file in the host contianer to our contianer
ADD . /go/src/github.com/daria/PortMicroserviceClient
WORKDIR /go/src/github.com/daria/PortMicroserviceClient
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN go get -u google.golang.org/grpc
RUN go install -v ./...
# Generate binary file from our /app
RUN go build /go/src/github.com/daria/PortMicroserviceClient/cmd/main.go
# Expose the ports used in server
EXPOSE 9080:9080
EXPOSE 3000:3000
# Run the app binarry file
CMD ["./main"]