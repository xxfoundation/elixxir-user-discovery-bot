module gitlab.com/elixxir/user-discovery-bot

go 1.13

require (
	github.com/golang/protobuf v1.4.3
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1
	gitlab.com/elixxir/client v1.2.1-0.20210112181113-e73779727e4c
	gitlab.com/elixxir/comms v0.0.4-0.20210108215809-3fd92afa101f
	gitlab.com/elixxir/crypto v0.0.7-0.20210107184400-5c3e52a35758
	gitlab.com/elixxir/primitives v0.0.3-0.20210107183456-9cf6fe2de1e5
	gitlab.com/xx_network/comms v0.0.4-0.20210107184305-7bb439c08ded
	gitlab.com/xx_network/crypto v0.0.5-0.20210107183440-804e0f8b7d22
	gitlab.com/xx_network/primitives v0.0.4-0.20210106014326-691ebfca3b07
	golang.org/x/sys v0.0.0-20210112091331-59c308dcf3cc // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210111234610-22ae2b108f89 // indirect
	google.golang.org/grpc v1.34.1 // indirect
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.27.1
