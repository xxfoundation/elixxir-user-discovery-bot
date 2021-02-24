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
	gitlab.com/elixxir/client v1.2.1-0.20210224024412-737e756de5d6
	gitlab.com/elixxir/comms v0.0.4-0.20210224021944-cd529dde0d9c
	gitlab.com/elixxir/crypto v0.0.7-0.20210216174551-f806f79610eb
	gitlab.com/elixxir/primitives v0.0.3-0.20210223210226-cccb5f7d4839
	gitlab.com/xx_network/comms v0.0.4-0.20210224021132-c5b1caeb0e58
	gitlab.com/xx_network/crypto v0.0.5-0.20210216174356-e81e1ddf8fb7
	gitlab.com/xx_network/primitives v0.0.4-0.20210219231511-983054dbee36
	golang.org/x/crypto v0.0.0-20210218145215-b8e89b74b9df // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210218155724-8ebf48af031b // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210218151259-fe80b386bf06 // indirect
	google.golang.org/grpc v1.35.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
