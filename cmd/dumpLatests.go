/*
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
package cmd

import (
	"fmt"
	"github.com/pestanko/isstat/app"
	"github.com/spf13/cobra"
	"os"
)

// dumpLatestsCmd represents the dumpLatests command
var dumpLatestsCmd = &cobra.Command{
	Use:   "dump-latest",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := app.GetAppConfig()
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(1)
		}

		application, err := app.GetApplication(&config)
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(1)
		}

		_, _ = application.DumpLatest()
	},
}

func init() {
	rootCmd.AddCommand(dumpLatestsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dumpLatestsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dumpLatestsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
