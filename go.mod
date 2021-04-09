module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.68 // indirect
	github.com/pelletier/go-toml v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/afero v1.5.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.5.1-0.20210409164625-0c3e924f95f6
	gitlab.com/elixxir/comms v0.0.4-0.20210409164117-ce218393248a
	gitlab.com/elixxir/crypto v0.0.7-0.20210409164001-5cb03c4d99b4
	gitlab.com/elixxir/primitives v0.0.3-0.20210409163616-5e22bc27ea48
	gitlab.com/xx_network/comms v0.0.4-0.20210406210737-45d1e87d294a
	gitlab.com/xx_network/crypto v0.0.5-0.20210405224157-2b1f387b42c1
	gitlab.com/xx_network/primitives v0.0.4-0.20210402222416-37c1c4d3fac4
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/genproto v0.0.0-20210406143921-e86de6bf7a46 // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
