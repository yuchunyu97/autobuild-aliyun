# autobuild-aliyun
基于阿里云容器镜像服务自动构建海外镜像工具。

使用方法
===

```bash
$ go get -u github.com/yuchunyu97/autobuild-aliyun/cmd/iproxy
$ iproxy
Use Alibaba Cloud ACR Service and git 
quickly pull foreign images. For example:

iproxy pull gcr.io/knative-releases/knative.dev/serving/cmd/queue:v0.14.0

Usage:
  iproxy [command]

Available Commands:
  help        Help about any command
  init        Initialize authentication information.
  login       Obtain temporary credentials for pulling images.
  pull        Pull an image.

Flags:
  -h, --help   help for iproxy

Use "iproxy [command] --help" for more information about a command.
                                                                               
```