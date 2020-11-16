module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/pkg/errors v0.9.1
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/client v1.2.1-0.20201116174929-8121d767912b
	gitlab.com/elixxir/comms v0.0.4-0.20201111191043-cce6aafab33b
	gitlab.com/elixxir/crypto v0.0.5-0.20201110193609-6b5e881867b4
	gitlab.com/xx_network/comms v0.0.4-0.20201110022115-4a6171cad07d
	gitlab.com/xx_network/crypto v0.0.4
	gitlab.com/xx_network/primitives v0.0.2
	golang.org/x/net v0.0.0-20201031054903-ff519b6c9102 // indirect
	golang.org/x/sys v0.0.0-20201101102859-da207088b7d1 // indirect
	google.golang.org/genproto v0.0.0-20201103154000-415bd0cd5df6 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
