module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.5.1-0.20210429183209-7dd569c43146
	gitlab.com/elixxir/comms v0.0.4-0.20210429182303-0edbda5e1b2c
	gitlab.com/elixxir/crypto v0.0.7-0.20210429182057-898c75770293
	gitlab.com/elixxir/primitives v0.0.3-0.20210429180244-cdbb97da0c16
	gitlab.com/xx_network/comms v0.0.4-0.20210426213447-82674e09e402
	gitlab.com/xx_network/crypto v0.0.5-0.20210420170153-2a6276844076
	gitlab.com/xx_network/primitives v0.0.4-0.20210402222416-37c1c4d3fac4
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.0.0-20210423184538-5f58ad60dda6 // indirect
	google.golang.org/genproto v0.0.0-20210427215850-f767ed18ee4d // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
