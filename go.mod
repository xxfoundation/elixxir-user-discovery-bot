module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.66 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.2.1-0.20210225221901-022ec8f3c25f
	gitlab.com/elixxir/comms v0.0.4-0.20210225221619-5a3e81a0f21f
	gitlab.com/elixxir/crypto v0.0.7-0.20210225184707-8e497d2c904e
	gitlab.com/elixxir/primitives v0.0.3-0.20210225184649-54d1b20caf89
	gitlab.com/xx_network/comms v0.0.4-0.20210225184643-04d57ac38237
	gitlab.com/xx_network/crypto v0.0.5-0.20210225184630-793a5fc60d3a
	gitlab.com/xx_network/primitives v0.0.4-0.20210225002641-4e446b2531ea
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7 // indirect
	golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210225212918-ad91960f0274 // indirect
	google.golang.org/grpc v1.36.0 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
