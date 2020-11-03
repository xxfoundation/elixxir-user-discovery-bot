module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/go-pg/pg v8.0.6+incompatible // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/mattn/go-shellwords v1.0.10 // indirect
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.10.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/assertions v1.1.0 // indirect
	github.com/spf13/afero v1.4.1 // indirect
	github.com/spf13/cobra v1.1.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0
	gitlab.com/elixxir/client v1.2.1-0.20201103000621-83d5c7af6e1d
	gitlab.com/elixxir/comms v0.0.3
	gitlab.com/elixxir/crypto v0.0.4
	gitlab.com/xx_network/comms v0.0.3
	gitlab.com/xx_network/crypto v0.0.4
	gitlab.com/xx_network/primitives v0.0.2
	golang.org/x/net v0.0.0-20201031054903-ff519b6c9102 // indirect
	golang.org/x/sys v0.0.0-20201101102859-da207088b7d1 // indirect
	google.golang.org/genproto v0.0.0-20201103154000-415bd0cd5df6 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
