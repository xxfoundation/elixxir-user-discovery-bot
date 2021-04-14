module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.5.2
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pelletier/go-toml v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.5.1-0.20210414202906-ed18541e447b
	gitlab.com/elixxir/comms v0.0.4-0.20210414200820-10e888270d4d
	gitlab.com/elixxir/crypto v0.0.7-0.20210412231025-6f75c577f803
	gitlab.com/elixxir/ekv v0.1.5 // indirect
	gitlab.com/elixxir/primitives v0.0.3-0.20210409190923-7bf3cd8d97e7
	gitlab.com/xx_network/comms v0.0.4-0.20210414191603-0904bc6eeda2
	gitlab.com/xx_network/crypto v0.0.5-0.20210413200952-56bd15ec9d99
	gitlab.com/xx_network/primitives v0.0.4-0.20210412170941-7ef69bce5a5c
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1 // indirect
	golang.org/x/sys v0.0.0-20210412220455-f1c623a9e750 // indirect
	google.golang.org/genproto v0.0.0-20210413151531-c14fb6ef47c3 // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
