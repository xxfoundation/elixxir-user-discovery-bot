package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gitlab.com/elixxir/client/storage/user"
	"gitlab.com/elixxir/client/xxdk"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/user-discovery-bot/banned"
	"gitlab.com/elixxir/user-discovery-bot/cmix"
	"gitlab.com/elixxir/user-discovery-bot/io"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/tls"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/ndf"
	"gitlab.com/xx_network/primitives/utils"
	"os"
	"time"
	"github.com/pkg/profile"

)

var (
	cfgFile, logPath                string
	certPath, keyPath, permCertPath string
	permAddress                     string
	logLevel                        uint // 0 = info, 1 = debug, >1 = trace
	validConfig                     bool
	devMode                         bool
	sessionPass                     string
)

const (
	profileMemFlag   = "profile-mem"
)

// RootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	Use:   "UDB",
	Short: "Runs the cMix UDB server.",
	Long:  "The cMix UDB server handles user and fact registration for the network.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		initLog()
		p := InitParams(viper.GetViper())

		memProfileOut := viper.GetString(profileMemFlag)
		if memProfileOut != "" {
			defer profile.Start(profile.MemProfile,
				profile.ProfilePath(memProfileOut),
				profile.NoShutdownHook).Stop()
		}


		// Initialize storage
		storage, err := storage.NewStorage(p.DbUsername, p.DbPassword, p.DbName, p.DbAddress, p.DbPort)
		if err != nil {
			jww.FATAL.Panicf("Failed to initialize storage interface: %+v", err)
		}

		_ = storage.StartFactManager(15 * time.Minute)

		var twilioManager *twilio.Manager
		devMode = viper.GetBool("devMode")
		if devMode {
			jww.WARN.Println("Twilio not configured; running with mock configuration")
			twilioManager = twilio.NewMockManager(storage)
		} else {
			twilioManager = twilio.NewManager(p.Twilio, storage)
		}

		cert, err := tls.LoadCertificate(string(p.PermCert))
		if err != nil {
			jww.FATAL.Fatalf("Failed to load permissioning cert to pem: %+v", err)
		}
		permCert, err := tls.ExtractPublicKey(cert)

		bannedManager, err := banned.NewManager(p.BannedUserList, p.BannedRegexList)
		if err != nil {
			jww.FATAL.Panicf("Failed to construct ban manager: %v", err)
		}

		jww.WARN.Printf("Skipping scheduling signature verification set to %v ", viper.GetBool("skipVerification"))

		// Set up manager with the ability to contact permissioning
		manager := io.NewManager(p.IO, &id.UDB, permCert, twilioManager,
			bannedManager, storage, viper.GetBool("skipVerification"))
		hostParams := connect.GetDefaultHostParams()
		hostParams.AuthEnabled = false
		permHost, err := manager.Comms.AddHost(&id.Permissioning,
			viper.GetString("permAddress"), p.PermCert, hostParams)
		if err != nil {
			jww.FATAL.Panicf("Unable to add permissioning host: %+v", err)
		}

		// Obtain the NDF from permissioning
		var returnedNdf *mixmessages.NDF
		retryTimer := 1 * time.Second
		for {
			returnedNdf, err = manager.Comms.RequestNdf(permHost)
			if err != nil {
				jww.WARN.Printf("Failed to get an ndf, Retying now: %s", err.Error())
				time.Sleep(retryTimer)
				continue
			}

			// Attempt to parse returned ndf
			parsedNdf, err := ndf.Unmarshal(returnedNdf.GetNdf())
			if err != nil {
				jww.ERROR.Printf("Unable to unmarshal returned NDF: %+v", err)
				time.Sleep(retryTimer)
				continue
			}

			// Verify the NDF has been populated
			if len(parsedNdf.Gateways) == 0 {
				jww.WARN.Printf("NDF is empty, waiting for gateways...")
				time.Sleep(retryTimer)
				continue
			}
			break
		}

		var user *xxdk.E2e
		cMixParams := xxdk.GetDefaultCMixParams()
		cMixParams.Network = cMixParams.Network.SetRealtimeOnlyAll()
		e2eParams := xxdk.GetDefaultE2EParams()
		if p.SessionPath != "" && utils.Exists(p.SessionPath) {
			// Construct a user using the NDF as a base
			user, err = LoginWithNDF(p.SessionPath, []byte(sessionPass),
				string(returnedNdf.GetNdf()),
				cMixParams, e2eParams)
			if err != nil {
				jww.FATAL.Fatalf("Failed to create user: %+v", err)
			}
		} else {
			user, err = LoginWithProto(p.SessionPath, []byte(sessionPass),
				p.ProtoUserJson, string(returnedNdf.GetNdf()),
				cMixParams, e2eParams)
			if err != nil {
				jww.FATAL.Fatalf("Failed to create user: %+v", err)
			}
		}

		m := cmix.NewManager(user, storage)
		m.Start()

		err = user.StartNetworkFollower(5 * time.Second)
		if err != nil {
			jww.FATAL.Fatal(err)
		}

		if err != nil {
			jww.FATAL.Panicf("%v", err)
		}
		// Wait forever

		select {}
	},
}

// LoginWithProto is a login function which constructs an xxdk.E2e object
// using a user.Proto which has been JSON marshalled.
func LoginWithProto(statePath string, statePass []byte,
	protoJson []byte, baseNdf string,
	cmixParams xxdk.CMIXParams, e2eParams xxdk.E2EParams) (*xxdk.E2e, error) {

	// Unmarshal the proto user JSON
	protoUser := &user.Proto{}
	err := json.Unmarshal(protoJson, protoUser)
	if err != nil {
		return nil, err
	}

	// Construct a network object
	err = xxdk.NewProtoCmix_Unsafe(baseNdf, statePath,
		statePass, protoUser)
	net, err := xxdk.LoadCmix(statePath,
		statePass, cmixParams)
	if err != nil {
		jww.FATAL.Panicf("%+v", err)
	}

	// Store the updated base NDF
	def, err := xxdk.ParseNDF(baseNdf)
	if err != nil {
		return nil, err
	}
	net.GetStorage().SetNDF(def)

	// Create a legacy identity
	identity, err := xxdk.MakeLegacyReceptionIdentity(net)
	if err != nil {
		return nil, err
	}

	// Create and return a messenger
	return xxdk.Login(net, xxdk.DefaultAuthCallbacks{}, identity, e2eParams)
}

// LoginWithNDF is a login function which creates an
// xxdk.E2e instance using a base NDF.
func LoginWithNDF(statePath string, statePass []byte, baseNdf string,
	cmixParams xxdk.CMIXParams, e2eParams xxdk.E2EParams) (*xxdk.E2e, error) {
	jww.INFO.Printf("LoginWithNDF()")

	// Construct a network object
	net, err := xxdk.LoadCmix(statePath, statePass, cmixParams)
	if err != nil {
		return nil, err
	}

	// Store the base NDF
	def, err := xxdk.ParseNDF(baseNdf)
	if err != nil {
		return nil, err
	}
	net.GetStorage().SetNDF(def)

	// Create a legacy identity
	identity, err := xxdk.MakeLegacyReceptionIdentity(net)
	if err != nil {
		return nil, err
	}

	// Create and return a messenger
	return xxdk.Login(net, xxdk.DefaultAuthCallbacks{}, identity, e2eParams)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		jww.ERROR.Println(err)
		os.Exit(1)
	}
}

func init() {
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
		"Path to the TLS certificate for UDB. Expects PEM format. Required field.")
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

	rootCmd.Flags().StringVar(&permAddress, "permAddress", "",
		"Public address of the Permissioning server. Required field.")
	err = viper.BindPFlag("permCertPath", rootCmd.Flags().Lookup("permCertPath"))
	handleBindingError(err, "permCertPath")

	rootCmd.Flags().StringVar(&sessionPass, "sessionPass", "", "Password for session files")
	err = viper.BindPFlag("sessionPass", rootCmd.Flags().Lookup("sessionPass"))
	handleBindingError(err, "sessionPass")

	rootCmd.Flags().BoolVarP(&devMode, "devMode", "", false, "Developer run mode")
	err = viper.BindPFlag("devMode", rootCmd.Flags().Lookup("devMode"))
	handleBindingError(err, "devMode")

	rootCmd.Flags().StringP("protoUserPath", "", "",
		"Path for ProtoUser file containing user primitives")
	err = viper.BindPFlag("protoUserPath", rootCmd.Flags().Lookup("protoUserPath"))
	handleBindingError(err, "protoUserPath")

	rootCmd.Flags().Bool("skipVerification", true,
		"Determines whether UD will verify a client's network signature "+
			"when registering. The default behaviour is to check the signature.")
	err = viper.BindPFlag("skipVerification", rootCmd.Flags().Lookup("skipVerification"))
	handleBindingError(err, "skipVerification")

	rootCmd.Flags().String(profileMemFlag, "",
		"Enable memory profiling to this directory")
	viper.BindPFlag(profileMemFlag, rootCmd.Flags().Lookup(profileMemFlag))

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
	var err error
	if cfgFile == "" {
		cfgFile, err = utils.SearchDefaultLocations("udb.yaml", "xxnetwork")
		if err != nil {
			validConfig = false
			jww.FATAL.Panicf("Failed to find config file: %+v", err)
		}
	} else {
		cfgFile, err = utils.ExpandPath(cfgFile)
		if err != nil {
			validConfig = false
			jww.FATAL.Panicf("Failed to expand config file path: %+v", err)
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
