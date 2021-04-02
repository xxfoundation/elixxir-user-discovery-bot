module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.68 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.5.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.5.1-0.20210402155653-ef72f4a44919
	gitlab.com/elixxir/comms v0.0.4-0.20210401210158-6053ad2e224c
	gitlab.com/elixxir/crypto v0.0.7-0.20210401210040-b7f1da24ef13
	gitlab.com/elixxir/primitives v0.0.3-0.20210401175645-9b7b92f74ec4
	gitlab.com/xx_network/comms v0.0.4-0.20210401160731-7b8890cdd8ad
	gitlab.com/xx_network/crypto v0.0.5-0.20210401160648-4f06cace9123
	gitlab.com/xx_network/primitives v0.0.4-0.20210331161816-ed23858bdb93
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210331212208-0fccb6fa2b5c // indirect
	golang.org/x/sys v0.0.0-20210402192133-700132347e07 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20210402141018-6c239bbf2bb1 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
