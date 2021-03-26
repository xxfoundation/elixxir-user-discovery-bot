module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.5.1
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.68 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.5.1-0.20210326172446-e613cf875b8d
	gitlab.com/elixxir/comms v0.0.4-0.20210326171912-e70c1821bf11
	gitlab.com/elixxir/crypto v0.0.7-0.20210326171146-c137bd7b0c6e
	gitlab.com/elixxir/primitives v0.0.3-0.20210326022836-1143187bd2fe
	gitlab.com/xx_network/comms v0.0.4-0.20210326005744-5e73cbf0f525
	gitlab.com/xx_network/crypto v0.0.5-0.20210319231335-249c6b1aa323
	gitlab.com/xx_network/primitives v0.0.4-0.20210309173740-eb8cd411334a
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210326060303-6b1517762897 // indirect
	google.golang.org/genproto v0.0.0-20210325224202-eed09b1b5210 // indirect
	google.golang.org/grpc v1.36.1 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
