module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.2.1-0.20210106014845-a7532b78dfa1
	gitlab.com/elixxir/comms v0.0.4-0.20210104224025-4e1cedc524a1
	gitlab.com/elixxir/crypto v0.0.7-0.20210104223925-7dfd3ad55d5c
	gitlab.com/elixxir/primitives v0.0.3-0.20210106014507-bf3dfe228fa6
	gitlab.com/xx_network/comms v0.0.4-0.20210106014446-be163ef3ccce
	gitlab.com/xx_network/crypto v0.0.5-0.20210106014410-0554a33a7124
	gitlab.com/xx_network/primitives v0.0.4-0.20210106014326-691ebfca3b07
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
