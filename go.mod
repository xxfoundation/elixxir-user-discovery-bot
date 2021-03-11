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
	gitlab.com/elixxir/client v1.5.1-0.20210311205740-89e522db8e79
	gitlab.com/elixxir/comms v0.0.4-0.20210309195247-fc17eb8560cf
	gitlab.com/elixxir/crypto v0.0.7-0.20210309193114-8a6225c667e2
	gitlab.com/elixxir/primitives v0.0.3-0.20210309193003-ef42ebb4800b
	gitlab.com/xx_network/comms v0.0.4-0.20210309192940-6b7fb39b4d01
	gitlab.com/xx_network/crypto v0.0.5-0.20210309192854-cf32117afb96
	gitlab.com/xx_network/primitives v0.0.4-0.20210309173740-eb8cd411334a
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210311153111-e2979279ddde // indirect
	google.golang.org/grpc v1.36.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
