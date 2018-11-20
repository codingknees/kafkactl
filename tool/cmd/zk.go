// Copyright © 2018 NAME HERE <jbonds@jbvm.io>
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
	"log"

	"github.com/spf13/cobra"
)

var (
	zkCommandInvoked bool
	zkTargetPath     string
	zkDeletePath     string
	zkTargetValue    string
	zkForceUpdate    bool
)

var zkCmd = &cobra.Command{
	Use:   "zk",
	Short: "Perform Various Zookeeper Administration Tasks",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zkCommandInvoked = true
		launchZKClient()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("create") {
			if cmd.Flags().Changed("value") {
				if zkTargetValue == "" {
					log.Fatalf("Error: Empty Value Submitted.\n")
				}
				targetVal := []byte(zkTargetValue)
				zkCreateValue(zkTargetPath, targetVal)
				return
			}
			zkCreateValue(zkTargetPath, nil)
			return
		}
		if cmd.Flags().Changed("delete") {
			zkDeleteValue(zkDeletePath)
			return
		}
		lsCmd.Run(cmd, args)
		return
	},
}

func init() {
	rootCmd.AddCommand(zkCmd)
	zkCmd.Flags().StringVarP(&zkTargetPath, "create", "c", "", "Create a Zookeeper Path (Use with --value for setting a value)")
	zkCmd.Flags().StringVarP(&zkDeletePath, "delete", "d", "", "Delete a Zookeeper Path/Value")
	zkCmd.Flags().StringVar(&zkTargetValue, "value", "", "Create a Zookeeper Value (Use with --create to specify the path for the value)")
	zkCmd.Flags().BoolVarP(&zkForceUpdate, "force", "f", false, "Force Operation")
}