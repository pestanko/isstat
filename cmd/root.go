package cmd

/**
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

import (
  "fmt"
  "github.com/pestanko/isstat/app"
  log "github.com/sirupsen/logrus"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "os"
)


var (
  cfgFile string
)


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "isstat",
  Short: "Get kontr statistics from the IS MUNI NOTEPADS",
  Long: `Get kontr statistics from the https://is.muni.cz notepads
         and convert them into JSON or CSV in order to analyze them.
        `,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  //	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {

  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/isstat/config.yaml)")
  rootCmd.PersistentFlags().StringP( "url", "U", "", "is muni url")
  rootCmd.PersistentFlags().StringP( "token", "T", "", "is muni token")
  rootCmd.PersistentFlags().StringP( "course", "C", "", "is muni course code")
  rootCmd.PersistentFlags().Int( "faculty-id", 0, "is muni faculty id")
  rootCmd.PersistentFlags().String( "parser", "", "parser for parsing the muni notepad content")
  rootCmd.PersistentFlags().String( "results", "", "results directory (default $CWD)")
  rootCmd.PersistentFlags().Bool( "dry-run", false, "dry run - do not execute the request")
  rootCmd.PersistentFlags().Bool( "without-timestamp", false, "create also without timestamp")

  _ = viper.BindPFlag("muni.url", rootCmd.PersistentFlags().Lookup("url"))
  _ = viper.BindPFlag("muni.token", rootCmd.PersistentFlags().Lookup("token"))
  _ = viper.BindPFlag("muni.course", rootCmd.PersistentFlags().Lookup("course"))
  _ = viper.BindPFlag("muni.faculty", rootCmd.PersistentFlags().Lookup("faculty-id"))
  _ = viper.BindPFlag("results", rootCmd.PersistentFlags().Lookup("results"))
  _ = viper.BindPFlag("parser", rootCmd.PersistentFlags().Lookup("parser"))
  _ = viper.BindPFlag("dryrun", rootCmd.PersistentFlags().Lookup("dry-run"))
  _ = viper.BindPFlag("without_timestamp", rootCmd.PersistentFlags().Lookup("without-timestamp"))

}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  app.SetupLogger("")

  if err := app.LoadConfig(cfgFile); err != nil {
    log.WithError(err).Error("Unable to load a config")
    os.Exit(1)
  }

}

