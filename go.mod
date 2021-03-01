module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.66 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.2.1-0.20210301192239-31caf4d6ee62
	gitlab.com/elixxir/comms v0.0.4-0.20210301173501-38cf2a1fc999
	gitlab.com/elixxir/crypto v0.0.7-0.20210226194937-5d641d5a31bc
	gitlab.com/elixxir/primitives v0.0.3-0.20210301183103-afcb5c703289
	gitlab.com/xx_network/comms v0.0.4-0.20210226194929-ea05928f74b7
	gitlab.com/xx_network/crypto v0.0.5-0.20210226194923-5f470e2a2533
	gitlab.com/xx_network/primitives v0.0.4-0.20210301172945-82f5d4248c04
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210226181700-f36f78243c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210226172003-ab064af71705 // indirect
	google.golang.org/grpc v1.36.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
