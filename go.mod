module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.67 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	github.com/ugorji/go v1.1.4 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
	gitlab.com/elixxir/client v1.5.1-0.20210316213105-df88c97e8852
	gitlab.com/elixxir/comms v0.0.4-0.20210316212548-21282c4ffdeb
	gitlab.com/elixxir/crypto v0.0.7-0.20210316212419-025ba777a80b
	gitlab.com/elixxir/primitives v0.0.3-0.20210309193003-ef42ebb4800b
	gitlab.com/xx_network/comms v0.0.4-0.20210316212314-d9a2a5e589c8
	gitlab.com/xx_network/crypto v0.0.5-0.20210316212152-504c0a1e0da0
	gitlab.com/xx_network/primitives v0.0.4-0.20210309173740-eb8cd411334a
	golang.org/x/crypto v0.0.0-20210314154223-e6e6c4f2bb5b // indirect
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	golang.org/x/sys v0.0.0-20210316164454-77fc1eacc6aa // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210315173758-2651cd453018 // indirect
	google.golang.org/grpc v1.36.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
