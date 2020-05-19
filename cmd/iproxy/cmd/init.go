/*Package cmd cmd
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
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yuchunyu97/autobuild-aliyun/pkg/auth"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize authentication information.",
	Long: `
Initialize Alibaba Cloud accessKey authentication information (only need to be done once).
`,
	Run: func(cmd *cobra.Command, args []string) {
		cred := auth.Credential{}

		fmt.Printf("Alibaba Cloud accessKeyID: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		cred.AccessKeyID = input.Text()

		fmt.Printf("Alibaba Cloud accessKeySecret: ")
		input = bufio.NewScanner(os.Stdin)
		input.Scan()
		cred.AccessKeySecret = input.Text()

		fmt.Printf("regionID: ")
		input = bufio.NewScanner(os.Stdin)
		input.Scan()
		cred.RegionID = input.Text()

		if err := cred.Save(); err != nil {
			fmt.Println("\nFail", err)
			return
		}

		fmt.Println("\nSucceed.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
