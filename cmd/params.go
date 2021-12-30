package cmd

import (
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/xx_network/primitives/utils"
	"net"
)

func InitParams(vip *viper.Viper) params.General {
	if !validConfig {
		jww.FATAL.Panicf("Invalid Config File: %s", cfgFile)
	}

	certPath = viper.GetString("certPath")
	if certPath == "" {
		jww.FATAL.Fatalf("certPath is blank - cannot run without certs")
	}
	cert, err := utils.ReadFile(certPath)
	if err != nil {
		jww.FATAL.Fatalf("Failed to read certificate at %s: %+v", certPath, err)
	}

	keyPath = viper.GetString("keyPath")
	if keyPath == "" {
		jww.FATAL.Fatalf("keyPath is blank - cannot run without keys")
	}
	key, err := utils.ReadFile(keyPath)
	if err != nil {
		jww.FATAL.Fatalf("Failed to read key at %s: %+v", keyPath, err)
	}

	permCertPath = viper.GetString("permCertPath")
	if permCertPath == "" {
		jww.FATAL.Fatalf("permCertPath is blank - cannot run without permissioning certificate")
	}
	permCert, err := utils.ReadFile(permCertPath)
	if err != nil {
		jww.FATAL.Fatalf("Failed to read permissioning certificate at %s: %+v", permCertPath, err)
	}

	sessionPath, err := utils.ExpandPath(viper.GetString("sessionPath"))
	if err != nil {
		jww.FATAL.Fatalf("Failed to read session path: %+v", err)
	}

	protoUserPath := viper.GetString("protoUserPath")
	if protoUserPath == "" {
		jww.FATAL.Fatalf("protoUserPath is blank - cannot run without proto user")
	}
	protoUserJson, err := utils.ReadFile(protoUserPath)
	if err != nil {
		jww.FATAL.Fatalf("Failed to read proto user at %s: %+v", protoUserPath, err)
	}

	sessionPass = viper.GetString("sessionPass")

	ioparams := params.IO{
		Cert: cert,
		Key:  key,
		Port: viper.GetString("port"),
	}

	// Obtain database connection info
	rawAddr := viper.GetString("dbAddress")
	var addr, port string
	if rawAddr != "" {
		addr, port, err = net.SplitHostPort(rawAddr)
		if err != nil {
			jww.FATAL.Panicf("Unable to get database port from %s: %+v", rawAddr, err)
		}
	}
	dbparams := params.Database{
		DbUsername: viper.GetString("dbUsername"),
		DbPassword: viper.GetString("dbPassword"),
		DbName:     viper.GetString("dbName"),
		DbAddress:  addr,
		DbPort:     port,
	}

	twilioparams := params.Twilio{
		AccountSid:      viper.GetString("twilioSid"),
		AuthToken:       viper.GetString("twilioToken"),
		VerificationSid: viper.GetString("twilioVerification"),
	}

	jww.INFO.Printf("config: %+v", viper.ConfigFileUsed())
	jww.INFO.Printf("Params: \n %+v", vip.AllSettings())
	jww.INFO.Printf("UDB port: %s", ioparams.Port)

	return params.General{
		PermCert:          permCert,
		SessionPath:       sessionPath,
		Database:          dbparams,
		IO:                ioparams,
		Twilio:            twilioparams,
		ProtoUserJsonPath: protoUserPath,
		ProtoUserJson:     protoUserJson,
	}
}
