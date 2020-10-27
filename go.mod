module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/go-pg/pg v8.0.6+incompatible // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mattn/go-shellwords v1.0.10
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.10.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/assertions v1.1.0 // indirect
	github.com/spf13/afero v1.4.1 // indirect
	github.com/spf13/cobra v1.1.0
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.2.1-0.20201015230424-9c4a4419d2c4
	gitlab.com/elixxir/comms v0.0.0-20201021175349-3249d073f1d1
	gitlab.com/elixxir/crypto v0.0.0-20201002151041-c4ab8f8033dc
	gitlab.com/elixxir/primitives v0.0.0-20200930214918-50b3c2030f26
	gitlab.com/xx_network/comms v0.0.0-20200925191822-08c0799a24a6
	gitlab.com/xx_network/crypto v0.0.0-20200812183430-c77a5281c686
	gitlab.com/xx_network/primitives v0.0.0-20200812183720-516a65a4a9b2
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee // indirect
	golang.org/x/sys v0.0.0-20201015000850-e3ed0017c211 // indirect
	google.golang.org/genproto v0.0.0-20201015140912-32ed001d685c // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
