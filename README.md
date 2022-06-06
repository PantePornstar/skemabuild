# skema-tool
dev tools for skema.  
Generate stubs and service code automatically from a single protobuf file.  
Get gRPC service code in one shot!  

## Quick Start

1. Installation  
First, make sure you have golang installed.  
Then setup necessary protocol buffers related tools:  
```shell
curl -L https://raw.githubusercontent.com/skema-dev/skema-tool/main/script/setup.sh | bash
```

2. Get Github Device Auth
```
st auth
```
You'll be prompted to open github device auth URL and input the token. Then just press Enter to save the github token on your local device.  


3. Init and Publish Protobuf API  
Now run and test(The following example is using my own sample repo. Change with your own account and repo)
```shell
# create an initial protocol buffers file
st api init --package=pack1 --service=hello

# create stubs for grpc-go and openapi, using protobuf file Hello1.proto, and output to ./stub-test
st api create --go_option github.com/likezhang-public/newst/test2/com.test/grpc-go --input ./Hello1.proto -o ./stub-test

# after stubs are created, upload to github and set the version tag
st api publish --stub=./stub-test --url  https://github.com/likezhang-public/newst/test2 --version=v0.0.2

# verify the package is downloadable
go get github.com/likezhang-public/newst/test2/com.test/grpc-go@v0.0.2
```
You can also skip creating local stubs and publish to github directly:  
```shell
st api init --package=pack1 --service=hello

st api publish --proto=./hello.proto --url  https://github.com/likezhang-public/newst/test2 --version=v0.0.2

go get github.com/likezhang-public/newst/test2/com.test/grpc-go@v0.0.2
```

4. Generate gRPC Service Code  
Now it's the final step: generate runnable gRPC service code automatically!  
Just use the following example to specify the protobuf file you published online  
```shell
st service create --proto https://github.com/likezhang-public/newst/blob/main/test2/com.test/grpc-go/Hello1.proto
```
Then we'll have the following code layout in current path. (you can specify output path by using --output)
.
├── cmd  
│     └── hello1service  
│         ├── config  
│         │   └── grpc.yaml  
│         ├── main.go  
│         └── service  
│             └── service.go  
├── go.mod  
└── internal  
    └── README.md   

To test it out, run `go mod tidy` first to install depenencies, then go to `cmd/hello1service` and `go run .`  
```shell
go mod tidy
cd cmd/hello1service
go run .
```
You'll see the launching information of grpc service and the http endpoint gateway:  
```
2022-06-07T00:48:32.460+0800	INFO	logging/logging.go:96	using local config from ./config/grpc.yaml
2022-06-07T00:48:32.460+0800	INFO	logging/logging.go:66	service port	{"gprc": 9991, "http": 9992}
2022-06-07T00:48:32.461+0800	INFO	logging/logging.go:96	gateway path is set to /
2022-06-07T00:48:32.461+0800	INFO	logging/logging.go:66	logging initialized:	{"level": "debug", "encoding": "console"}
2022-06-07T00:48:32.461+0800	INFO	logging/logging.go:96	Serving gRPC start...
```
To further understand how the code is designed, Please checkout our [Skema-Go Framework](https://github.com/skema-dev/skema-go)
