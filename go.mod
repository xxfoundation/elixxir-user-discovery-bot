module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/aws/aws-lambda-go v1.8.1 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.68 // indirect
	github.com/pelletier/go-toml v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.5.1-0.20210427005915-c57544aa7d4d
	gitlab.com/elixxir/comms v0.0.4-0.20210427005410-7ae183abda2a
	gitlab.com/elixxir/crypto v0.0.7-0.20210427005255-4fe1bcf69c5a
	gitlab.com/elixxir/primitives v0.0.3-0.20210427004615-c68ecf15fcf3
	gitlab.com/xx_network/comms v0.0.4-0.20210426213447-82674e09e402
	gitlab.com/xx_network/crypto v0.0.5-0.20210420170153-2a6276844076
	gitlab.com/xx_network/primitives v0.0.4-0.20210402222416-37c1c4d3fac4
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.0.0-20210423184538-5f58ad60dda6 // indirect
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
	google.golang.org/genproto v0.0.0-20210426193834-eac7f76ac494 // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
