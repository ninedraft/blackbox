// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//var cfgFile string

type BoxMode int

const (
	ModeServe BoxMode = iota
	ModeRefine
	ModeInit
)

var Configuration struct {
	Addr           string
	Mode           BoxMode
	DBFilePath     string
	Proxy          string
	ConfigFile     string
	LogFormatParam string
	LogLevelParam  string
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "blackbox",
	Short: "A minimalistic blog engine",
	Long: `Blackbox is a minimalisic blog engine, 
which provides a simple building, configuration and runnig solution.
Powered by well trained gophers!`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().
		StringVar(&Configuration.ConfigFile, "config", "./blackbox.yaml", "config file")
	RootCmd.PersistentFlags().
		StringVarP(&Configuration.DBFilePath, "db", "", "./blackbox.db", "embedded db file")
	//initCmd.PersistentFlags().
	//	StringVar(&Configuration.NginxConfig, "proxyconfig", "proxy.config", "proxy config file")
	RootCmd.PersistentFlags().
		StringVarP(&Configuration.Proxy, "proxy", "p", "", "run proxy on frontend, [nginx|caddy]:\"path to custom config if need\"")
	RootCmd.PersistentFlags().
		StringVar(&Configuration.LogFormatParam, "logf", "text", "log output format: json | text")
	RootCmd.PersistentFlags().
		StringVar(&Configuration.LogLevelParam, "logl", "info", "log verbosity level: debug | info | fatal | panic")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if Configuration.ConfigFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(Configuration.ConfigFile)
	}

	viper.SetConfigName(".blackbox") // name of config file (without extension)
	viper.AddConfigPath("")
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
