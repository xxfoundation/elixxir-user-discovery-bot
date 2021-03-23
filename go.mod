module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.5.1
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.68 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.5.1-0.20210323170252-7aa7a34f2682
	gitlab.com/elixxir/comms v0.0.4-0.20210323165848-495b7efbf1af
	gitlab.com/elixxir/crypto v0.0.7-0.20210319231554-b73b6e62ddbc
	gitlab.com/elixxir/primitives v0.0.3-0.20210309193003-ef42ebb4800b
	gitlab.com/xx_network/comms v0.0.4-0.20210323140408-2b2613abb5a3
	gitlab.com/xx_network/crypto v0.0.5-0.20210319231335-249c6b1aa323
	gitlab.com/xx_network/primitives v0.0.4-0.20210309173740-eb8cd411334a
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210322173543-5f0e89347f5a // indirect
	google.golang.org/grpc v1.36.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
