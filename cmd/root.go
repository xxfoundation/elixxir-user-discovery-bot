////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Package cmd initializes the CLI and config parsers as well as the logger.
package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gitlab.com/elixxir/client/api"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/primitives/utils"
	"gitlab.com/elixxir/user-discovery-bot/cmix"
	"net"
	"os"
)

var cfgFile string
var logLevel uint // 0 = info, 1 = debug, >1 = trace
var noTLS bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "user-discovery-bot",
	Short: "Runs a user discovery bot for cMix",
	Long:  `This bot provides user lookup and search functions on cMix`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		sess := viper.GetString("sessionfile")
		if sess == "" {
			sess = "udb-session.blob"
		}

		cmix.BannedUsernameList = *cmix.InitBlackList(viper.GetString("blacklistedNamesFilePath"))

		// Set up database connection
		rawAddr := viper.GetString("dbAddress")

		var addr, port string
		var err error
		if rawAddr != "" {
			addr, port, err = net.SplitHostPort(rawAddr)
			if err != nil {
				jww.FATAL.Panicf("Unable to get database port: %+v", err)
			}
		}

		//params := Params{
		//	viper.GetString("dbUsername"),
		//	viper.GetString("dbPassword"),
		//	viper.GetString("dbName"),
		//	addr,
		//	port,
		//	viper.GetString("sessionfile"),
		//	viper.GetString("ndfPath"),
		//}
		//if params.sessionPath == "" {
		//	params.sessionPath = "udb-session.blob"
		//}

		// Import the network definition file
		ndfBytes, err := utils.ReadFile()
		if err != nil {
			globals.Log.FATAL.Panicf("Could not read network definition file: %v", err)
		}
		ndfJSON := api.VerifyNDF(string(ndfBytes), "")

		err = StartBot(sess, ndfJSON)
		if err != nil {
			globals.Log.FATAL.Panicf("Could not start bot: %+v", err)
		}
		// Block forever as a keepalive
		select {}
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately.  This is called by main.main(). It only needs to
// happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		cmix.Log.ERROR.Println(err)
		os.Exit(1)
	}
}

// init is the initialization function for Cobra which defines commands
// and flags.
func init() {
	cmix.Log.DEBUG.Print("Printing log from init")
	// NOTE: The point of init() is to be declarative.
	// There is one init in each sub command. Do not put variable declarations
	// here, and ensure all the Flags are of the *P variety, unless there's a
	// very good reason not to have them as local params to sub command."
	cobra.OnInitialize(initConfig, initLog)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.Flags().StringVarP(&cfgFile, "config", "", "",
		"config file (default is $PWD/udb.yaml)")
	RootCmd.Flags().UintVarP(&logLevel, "logLevel", "l", 1,
		"Level of debugging to display. 0 = info, 1 = debug, >1 = trace")
	RootCmd.Flags().BoolVarP(&noTLS, "noTLS", "", false,
		"Set to ignore TLS")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		// Default search paths
		var searchDirs []string
		searchDirs = append(searchDirs, "./") // $PWD
		// $HOME
		home, _ := homedir.Dir()
		searchDirs = append(searchDirs, home+"/.elixxir/")
		// /etc/elixxir
		searchDirs = append(searchDirs, "/etc/elixxir")
		jww.DEBUG.Printf("Configuration search directories: %v", searchDirs)

		for i := range searchDirs {
			cfgFile = searchDirs[i] + "udb.yaml"
			_, err := os.Stat(cfgFile)
			if !os.IsNotExist(err) {
				break
			}
		}
	}
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Unable to read config file (%s): %+v", cfgFile, err.Error())
	}

}

// initLog initializes logging thresholds and the log path.
func initLog() {
	vipLogLevel := viper.GetUint("logLevel")

	// Check the level of logs to display
	if vipLogLevel > 1 {
		// Set the GRPC log level
		if vipLogLevel > 1 {
			err := os.Setenv("GRPC_GO_LOG_SEVERITY_LEVEL", "info")
			if err != nil {
				jww.ERROR.Printf("Could not set GRPC_GO_LOG_SEVERITY_LEVEL: %+v", err)
			}

			err = os.Setenv("GRPC_GO_LOG_VERBOSITY_LEVEL", "99")
			if err != nil {
				jww.ERROR.Printf("Could not set GRPC_GO_LOG_VERBOSITY_LEVEL: %+v", err)
			}
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

	if viper.Get("logPath") != nil {
		// Create log file, overwrites if existing
		logPath := viper.GetString("logPath")
		logFile, err := os.Create(logPath)
		if err != nil {
			cmix.Log.WARN.Println("Invalid or missing log path, default path used.")
		} else {
			cmix.Log.SetLogOutput(logFile)
		}
	}
}
