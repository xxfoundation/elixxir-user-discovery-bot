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
	gitlab.com/elixxir/client v1.2.1-0.20210308200451-2d5972c95e9b
	gitlab.com/elixxir/comms v0.0.4-0.20210308190443-cc7781473650
	gitlab.com/elixxir/crypto v0.0.7-0.20210305221450-c94429c34886
	gitlab.com/elixxir/primitives v0.0.3-0.20210305221357-0b3f19cc3c8a
	gitlab.com/xx_network/comms v0.0.4-0.20210305221336-a7e28febdbf6
	gitlab.com/xx_network/crypto v0.0.5-0.20210305221255-f0ef174f5135
	gitlab.com/xx_network/primitives v0.0.4-0.20210308175329-436b0c8753ea
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210305215415-5cdee2b1b5a0 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210303154014-9728d6b83eeb // indirect
	google.golang.org/grpc v1.36.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
