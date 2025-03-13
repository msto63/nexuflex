module github.com/nexuflex/nexuflex-server

go 1.20

require (
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.3.1
	github.com/nexuflex/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.58.2
	google.golang.org/protobuf v1.31.0
)

require (
	golang.org/x/crypto v0.13.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

replace github.com/nexuflex/shared => ../shared 
