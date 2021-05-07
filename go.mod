module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	github.com/ugorji/go v1.1.4 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
	gitlab.com/elixxir/client v1.5.1-0.20210506194805-11bbfc1912c9
	gitlab.com/elixxir/comms v0.0.4-0.20210506200648-51f81a7e49a3
	gitlab.com/elixxir/crypto v0.0.7-0.20210506223047-3196e4301110
	gitlab.com/elixxir/primitives v0.0.3-0.20210504210415-34cf31c2816e
	gitlab.com/xx_network/comms v0.0.4-0.20210506193128-5af6bddf0ae0
	gitlab.com/xx_network/crypto v0.0.5-0.20210506192937-7882aa3810b4
	gitlab.com/xx_network/primitives v0.0.4-0.20210506192747-def158203920
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.0.0-20210423184538-5f58ad60dda6 // indirect
	google.golang.org/genproto v0.0.0-20210427215850-f767ed18ee4d // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
