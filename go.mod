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
	gitlab.com/elixxir/client v1.2.1-0.20210217224525-69f8857780b6
	gitlab.com/elixxir/comms v0.0.4-0.20210217173658-d751b9c00e5a
	gitlab.com/elixxir/crypto v0.0.7-0.20210216174551-f806f79610eb
	gitlab.com/elixxir/primitives v0.0.3-0.20210216174458-2a23825c1eb1
	gitlab.com/xx_network/comms v0.0.4-0.20210216174438-0790d1f1f225
	gitlab.com/xx_network/crypto v0.0.5-0.20210216174356-e81e1ddf8fb7
	gitlab.com/xx_network/primitives v0.0.4-0.20210215192713-e32335847d4f
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210217105451-b926d437f341 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210217220511-c18582744cc2 // indirect
	google.golang.org/grpc v1.35.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
