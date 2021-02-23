module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.66 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	github.com/ugorji/go v1.1.4 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
	gitlab.com/elixxir/client v1.2.1-0.20210223210950-8433deb3e64a
	gitlab.com/elixxir/comms v0.0.4-0.20210223210438-35bdbcaf78be
	gitlab.com/elixxir/crypto v0.0.7-0.20210223210315-b2072c080b0f
	gitlab.com/elixxir/primitives v0.0.3-0.20210223210226-cccb5f7d4839
	gitlab.com/xx_network/comms v0.0.4-0.20210223210205-6d1cb7fde5d1
	gitlab.com/xx_network/crypto v0.0.5-0.20210223210125-9c1a8a8f1ec6
	gitlab.com/xx_network/primitives v0.0.4-0.20210219231511-983054dbee36
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210222171744-9060382bd457 // indirect
	golang.org/x/sys v0.0.0-20210223095934-7937bea0104d // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210223151946-22b48be4551b // indirect
	google.golang.org/grpc v1.35.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
