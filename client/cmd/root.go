/*
 Copyright 2017 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/
package cmd

import (
	log "github.com/Sirupsen/logrus"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RED, GREEN func(a ...interface{}) string

var cfgFile string
var KubeconfigPath string
var Labelselector string
var DebugFlag bool
var Namespace string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pgo",
	Short: "The pgo command line interface.",
	Long: `The pgo command line interface lets you
create and manage PostgreSQL clusters.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Debug(err.Error())
		os.Exit(-1)
	}

}

func init() {

	GREEN = color.New(color.FgGreen).SprintFunc()
	RED = color.New(color.FgRed).SprintFunc()

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pgo.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().StringVar(&KubeconfigPath, "kubeconfig", "", "kube config file")
	RootCmd.PersistentFlags().StringVar(&Namespace, "namespace", "", "kube namespace to work in (default is default)")
	RootCmd.PersistentFlags().StringVar(&Labelselector, "selector", "", "label selector string")
	RootCmd.PersistentFlags().BoolVar(&DebugFlag, "debug", false, "enable debug with true")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".pgo")     // name of config file (without extension)
	viper.AddConfigPath(".")        // adding current directory as first search path
	viper.AddConfigPath("$HOME")    // adding home directory as second search path
	viper.AddConfigPath("/etc/pgo") // adding /etc/pgo directory as third search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		log.Debugf("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Debug("config file not found")
	}

	if DebugFlag || viper.GetBool("PGO.DEBUG") {
		log.Debug("debug flag is set to true")
		log.SetLevel(log.DebugLevel)
	}

	if KubeconfigPath == "" {
		KubeconfigPath = viper.GetString("KUBECONFIG")
	}
	if KubeconfigPath == "" {
		log.Error("--kubeconfig flag is not set and required")
		os.Exit(2)
	}

	log.Debug("kubeconfig path is " + viper.GetString("KUBECONFIG"))

	if Namespace == "" {
		Namespace = viper.GetString("NAMESPACE")
	}
	if Namespace == "" {
		log.Error("--namespace flag is not set and required")
		os.Exit(2)
	}

	log.Debug("namespace is " + viper.GetString("NAMESPACE"))
	ConnectToKube()

	/*
		file, err2 := os.Create("/tmp/pgo-bash-completion.out")
		if err2 != nil {
			log.Error(err2.Error())
		}
		defer file.Close()
		RootCmd.GenBashCompletion(file)
	*/

}
