module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.2.1-0.20210105235934-ee063885ac30
	gitlab.com/elixxir/comms v0.0.4-0.20210104224025-4e1cedc524a1
	gitlab.com/elixxir/crypto v0.0.7-0.20210104223925-7dfd3ad55d5c
	gitlab.com/elixxir/ekv v0.1.4 // indirect
	gitlab.com/elixxir/primitives v0.0.3-0.20210104223605-0e47af99d9d5
	gitlab.com/xx_network/comms v0.0.4-0.20201217200138-87075d5b4ffd
	gitlab.com/xx_network/crypto v0.0.5-0.20201217195719-cc31e1d1eee3
	gitlab.com/xx_network/primitives v0.0.4-0.20201216174909-808eb0fc97fc
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/sys v0.0.0-20210105210732-16f7687f5001 // indirect
	google.golang.org/genproto v0.0.0-20210105202744-fe13368bc0e1 // indirect
	google.golang.org/grpc v1.34.0 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
