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
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuchunyu97/autobuild-aliyun/pkg/acr"
	"github.com/yuchunyu97/autobuild-aliyun/pkg/auth"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Obtain temporary credentials for pulling images.",
	Long: `
Obtain temporary credentials for pulling images
`,
	Run: func(cmd *cobra.Command, args []string) {
		// 先读取认证信息
		cred := auth.Credential{}
		if err := cred.Get(); err != nil {
			fmt.Println(err)
			return
		}
		// 初始化一个阿里云客户端
		client, err := acr.NewClient(
			cred.RegionID,
			cred.AccessKeyID,
			cred.AccessKeySecret,
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		res, err := client.GetAuthorizationToken()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(res)
		return
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
