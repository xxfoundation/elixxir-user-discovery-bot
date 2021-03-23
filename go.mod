module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.5.1
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.68 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.5.1-0.20210323234532-721433cdf99e
	gitlab.com/elixxir/comms v0.0.4-0.20210323234032-57e81abb3171
	gitlab.com/elixxir/crypto v0.0.7-0.20210319231554-b73b6e62ddbc
	gitlab.com/elixxir/primitives v0.0.3-0.20210309193003-ef42ebb4800b
	gitlab.com/xx_network/comms v0.0.4-0.20210323233204-5acf90f56550
	gitlab.com/xx_network/crypto v0.0.5-0.20210319231335-249c6b1aa323
	gitlab.com/xx_network/primitives v0.0.4-0.20210309173740-eb8cd411334a
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210323141857-08027d57d8cf // indirect
	google.golang.org/genproto v0.0.0-20210323160006-e668133fea6a // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
