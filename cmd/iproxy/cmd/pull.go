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
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yuchunyu97/autobuild-aliyun/pkg/acr"
	"github.com/yuchunyu97/autobuild-aliyun/pkg/auth"
	"github.com/yuchunyu97/autobuild-aliyun/pkg/code"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull an image.",
	Long: `
Pull an image. For example:

iproxy pull gcr.io/knative-releases/knative.dev/serving/cmd/queue:v0.14.0`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 读取要拉取的镜像地址
		imageName := args[0]
		err := pullImage(imageName)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const (
	repoNamespace = "iproxy"
	repoName      = "autobuild"
)

func pullImage(imageName string) (err error) {
	imagePath := strings.Replace(imageName, ":", "/", -1)
	imagePath = strings.Replace(imagePath, "@", "/", -1)
	if len(imagePath) > 128 {
		imagePath = imagePath[:128] // imagePath 长度限制 128
	}
	m := md5.Sum([]byte(imageName))
	imageMD5 := fmt.Sprintf("%x", m)

	filename := "Dockerfile"
	content := fmt.Sprintf(`# Date: %s
# MD5: %s
FROM %s
	`, time.Now().Format("2006-01-02 15:04:05"), imageMD5, imageName)

	// =======更新代码=======

	// 下载代码库
	fmt.Printf("正在更新代码")
	repo, err := code.PrepareCode()
	if err != nil {
		return
	}
	// 添加文件
	err = repo.AddFile(
		imagePath,
		filename,
		content)
	if err != nil {
		return
	}
	// 提交 Git
	err = repo.Submit(imageName)
	if err != nil {
		return
	}
	fmt.Printf("\r更新代码成功\n")
	// 删除临时目录
	repo.Remove()

	// =======构建镜像=======

	// 先读取认证信息
	cred := auth.Credential{}
	if err = cred.Get(); err != nil {
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
		return
	}

	// 获取现有构建规则列表，并找到可用的构建规则 ID，更新构建规则
	var buildRuleID int
	repoBuildRules, err := client.GetRepoBuildRuleList(repoNamespace, repoName)
	if err != nil {
		return
	}
	if len(repoBuildRules.BuildRules) == 0 {
		// 新建构建规则
		repoNewBuildRules, err := client.CreateRepoBuildRule(repoNamespace, repoName, imagePath, imageMD5)
		if err != nil {
			return err
		}
		buildRuleID = repoNewBuildRules.BuildRuleID
	} else {
		// 更新构建规则
		repoReBuildRules, err := client.UpdateRepoBuildRule(repoNamespace, repoName, imagePath, imageMD5, repoBuildRules.BuildRules[0])
		if err != nil {
			return err
		}
		buildRuleID = repoReBuildRules.BuildRuleID
	}

	// 触发构建
	fmt.Println("准备拉取镜像")
	_, err = client.StartRepoBuildByRule(repoNamespace, repoName, buildRuleID)
	if err != nil {
		fmt.Println("拉取镜像失败")
		return err
	}

	// 通过列表获取构建的 buildID
	var buildID string
	resBuilds, err := client.GetRepoBuildList(repoNamespace, repoName)
	if err != nil {
		return err
	}
	// 数组按时间倒序
	length := len(resBuilds.Builds)
	for i := 0; i < length/2; i++ {
		temp := resBuilds.Builds[length-1-i]
		resBuilds.Builds[length-1-i] = resBuilds.Builds[i]
		resBuilds.Builds[i] = temp
	}
	for _, v := range resBuilds.Builds {
		if v.Image.Tag == imageMD5 {
			buildID = v.BuildID
		}
	}

	// 轮询是否构建完成
	for {
		resBuild, err := client.GetRepoBuildStatus(repoNamespace, repoName, buildID)
		if err != nil {
			return fmt.Errorf("轮询检查失败 : %s", err)
		}

		status := resBuild.BuildStatus
		fmt.Println(status)
		if status == "SUCCESS" {
			fmt.Println("镜像拉取成功")
			break
		} else if status == "PENDING" || status == "BUILDING" {
			time.Sleep(time.Second * 2)
			continue
		} else {
			// status == "FAILED" || status == "CANCELED"
			return fmt.Errorf("镜像拉取失败 Build ID: %s", buildID)
		}
	}

	// 输出拉取镜像的命令
	fmt.Printf(`
使用如下命令拉取镜像：

docker pull registry.cn-qingdao.aliyuncs.com/%[3]s/%[4]s:%[1]s
docker tag registry.cn-qingdao.aliyuncs.com/%[3]s/%[4]s:%[1]s %[2]s
docker rmi registry.cn-qingdao.aliyuncs.com/%[3]s/%[4]s:%[1]s

提示：请使用 iproxy login 来获取拉取镜像需要的临时凭证信息。
`, imageMD5, imageName, repoNamespace, repoName)
	return
}
