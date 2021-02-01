package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/user-discovery-bot/cmix"
	"gitlab.com/elixxir/user-discovery-bot/io"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/crypto/tls"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/utils"
	"os"
)

var (
	cfgFile, logPath                string
	certPath, keyPath, permCertPath string
	logLevel                        uint // 0 = info, 1 = debug, >1 = trace
	validConfig                     bool
	devMode                         bool
	sessionPass                     string
)

// RootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	Use:   "UDB",
	Short: "Runs the cmix UDB server",
	Long:  `The cMix UDB server handles user & fact registration for the network`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		initLog()
		p := InitParams(viper.GetViper())
		storage, _, err := storage.NewStorage(p.Database)

		var twilioManager *twilio.Manager
		if devMode {
			twilioManager = twilio.NewMockManager(storage)
		} else {
			twilioManager = twilio.NewManager(p.Twilio, storage)
		}

		cert, err := tls.LoadCertificate(string(p.PermCert))
		if err != nil {
			jww.FATAL.Fatalf("Failed to load permissioning cert to pem: %+v", err)
		}
		permCert, err := tls.ExtractPublicKey(cert)
		_ = io.NewManager(p.IO, &id.UDB, permCert, twilioManager, storage)

		m, err := cmix.NewManager(p.SessionPath, []byte(sessionPass), storage)
		if err != nil {
			jww.FATAL.Fatalf("Failed to create cmix manager: %+v", err)
		}

		err = m.Start()
		if err != nil {
			jww.FATAL.Fatalf("Failed to start cmix manager: %+v", err)
		}
		// Wait forever
		select {}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		jww.ERROR.Println(err)
		os.Exit(1)
	}
}

func init() {
	// NOTE: The point of init() is to be declarative.
	// There is one init in each sub command. Do not put variable declarations
	// here, and ensure all the Flags are of the *P variety, unless there's a
	// very good reason not to have them as local Params to sub command."

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "",
		"Path to load the UDB configuration file from. If not set, this "+
			"file must be named udb.yaml and must be located in "+
			"~/.xxnetwork/, /opt/xxnetwork, or /etc/xxnetwork.")

	rootCmd.Flags().IntP("port", "p", -1,
		"Port for UDB to listen on. UDB must be the only listener "+
			"on this port. Required field.")
	err := viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	handleBindingError(err, "port")

	rootCmd.Flags().UintVarP(&logLevel, "logLevel", "l", 0,
		"Level of debugging to print (0 = info, 1 = debug, >1 = trace).")
	err = viper.BindPFlag("logLevel", rootCmd.Flags().Lookup("logLevel"))
	handleBindingError(err, "logLevel")

	rootCmd.Flags().StringVar(&logPath, "log", "./udb-logs/udb.log",
		"Path where log file will be saved.")
	err = viper.BindPFlag("log", rootCmd.Flags().Lookup("log"))
	handleBindingError(err, "log")

	rootCmd.Flags().StringVar(&certPath, "certPath", "",
		"Path to the TLS certificate for UDB. Expects PEM "+
			"format. Required field.")
	err = viper.BindPFlag("certPath", rootCmd.Flags().Lookup("certPath"))
	handleBindingError(err, "certPath")

	rootCmd.Flags().StringVar(&keyPath, "keyPath", "",
		"Path to the private key associated with UDB TLS "+
			"certificate. Required field.")
	err = viper.BindPFlag("keyPath", rootCmd.Flags().Lookup("keyPath"))
	handleBindingError(err, "keyPath")

	rootCmd.Flags().StringVar(&permCertPath, "permCertPath", "",
		"Path to the TLS certificate for Permissioning server. Expects PEM "+
			"format. Required field.")
	err = viper.BindPFlag("permCertPath", rootCmd.Flags().Lookup("permCertPath"))
	handleBindingError(err, "permCertPath")

	rootCmd.Flags().StringVar(&sessionPass, "sessionPass", "", "Password for session files")
	err = viper.BindPFlag("sessionPass", rootCmd.Flags().Lookup("sessionPass"))
	handleBindingError(err, "sessionPass")

	rootCmd.Flags().BoolVarP(&devMode, "devMode", "", false, "Developer run mode")
	err = viper.BindPFlag("devMode", rootCmd.Flags().Lookup("devMode"))
	handleBindingError(err, "devMode")
}

// Handle flag binding errors
func handleBindingError(err error, flag string) {
	if err != nil {
		jww.FATAL.Panicf("Error on binding flag \"%s\":%+v", flag, err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	validConfig = true
	if cfgFile == "" {
		var err error
		cfgFile, err = utils.SearchDefaultLocations("udb.yaml", "xxnetwork")
		if err != nil {
			validConfig = false
			jww.FATAL.Panicf("Failed to find config file: %+v", err)
		}
	}
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Unable to read config file (%s): %+v", cfgFile, err.Error())
		validConfig = false
	}
}

// initLog initializes logging thresholds and the log path.
func initLog() {
	vipLogLevel := viper.GetUint("logLevel")

	// Check the level of logs to display
	if vipLogLevel > 1 {
		// Set the GRPC log level
		err := os.Setenv("GRPC_GO_LOG_SEVERITY_LEVEL", "info")
		if err != nil {
			jww.ERROR.Printf("Could not set GRPC_GO_LOG_SEVERITY_LEVEL: %+v", err)
		}

		err = os.Setenv("GRPC_GO_LOG_VERBOSITY_LEVEL", "99")
		if err != nil {
			jww.ERROR.Printf("Could not set GRPC_GO_LOG_VERBOSITY_LEVEL: %+v", err)
		}
		// Turn on trace logs
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelTrace)
		mixmessages.TraceMode()
	} else if vipLogLevel == 1 {
		// Turn on debugging logs
		jww.SetLogThreshold(jww.LevelDebug)
		jww.SetStdoutThreshold(jww.LevelDebug)
		mixmessages.DebugMode()
	} else {
		// Turn on info logs
		jww.SetLogThreshold(jww.LevelInfo)
		jww.SetStdoutThreshold(jww.LevelInfo)
	}

	logPath = viper.GetString("log")

	logFile, err := os.OpenFile(logPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644)
	if err != nil {
		fmt.Printf("Could not open log file %s!\n", logPath)
	} else {
		jww.SetLogOutput(logFile)
	}
}
