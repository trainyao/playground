# 使用 kubebuilder 重构原始方式的 CRD 和 Controller

## kubebuilder 初始化工程
[![asciicast](https://asciinema.org/a/475539.svg)](https://asciinema.org/a/475539?speed=2&i=0.3)

## kubebuilder main.go controller.go 逻辑迁移

这里 asciinema cast 文件过大了, 无法上传, 可以下载本目录的 `kubebuilder-controller.cast` 然后执行

```zsh

asciinema play kubebuilder-controller.cast -i 0.3 -s 2

```

观看回放

## 漏掉一个小细节, types.go 内需要去掉 status subresource 的注释, 否则会 patch status 失败
[![asciicast](https://asciinema.org/a/475540.svg)](https://asciinema.org/a/475540?speed=2&i=0.3)

## 部署
[![asciicast](https://asciinema.org/a/475541.svg)](https://asciinema.org/a/475541?speed=2&i=0.3)

