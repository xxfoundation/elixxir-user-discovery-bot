module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.61 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.2.1-0.20210125232058-a819eceb2f1d
	gitlab.com/elixxir/comms v0.0.4-0.20210125231706-c782d1f645ea
	gitlab.com/elixxir/crypto v0.0.7-0.20210125231257-b52a556be893
	gitlab.com/elixxir/primitives v0.0.3-0.20210125225949-9469ce6b08fc
	gitlab.com/xx_network/comms v0.0.4-0.20210121204701-7a1eb0542424
	gitlab.com/xx_network/crypto v0.0.5-0.20210121204626-b251b926e4f7
	gitlab.com/xx_network/primitives v0.0.4-0.20210121203635-8a771fc14f8a
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210125195502-f46fe6c6624a // indirect
	google.golang.org/grpc v1.35.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
