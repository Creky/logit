# 📝 logit

[![License](./license.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

**logit** 是一个简单易用并且是基于级别控制的日志库，可以应用于所有的 [GoLang](https://golang.org) 应用程序中。

[Read me in English](./README.en.md).

### 🥇 功能特性

* 支持日志级别控制，目前一共有四个日志级别
* 支持开启或者关闭日志功能，线上环境可以关闭或调高日志级别

### 🚀 安装方式

唯一需要的依赖就是 [Golang 运行环境](https://golang.org).

> Go modules

```bash
$ go get github.com/FishGoddess/logit@v0.0.2
```

您也可以直接编辑 go.mod 文件，然后执行 _**go build**_.

```bash
module your_project_name

go 1.14

require (
    github.com/FishGoddess/logit v0.0.2
)
```

> Go path

```bash
$ go get -u github.com/FishGoddess/logit
```

logit 没有任何其他额外的依赖，纯使用 [Golang 标准库](https://golang.org) 完成。

```go
package main

import (
    "github.com/FishGoddess/logit"
)

func main() {
    
    // log as you want.
    logit.Debug("I am a debug message! But I will not be logged in default level!")
    logit.Info("I am an info message!")
    logit.Warning("I am a warning message!")
    logit.Error("I am an error message!")
    
    // change log level.
    logit.ChangeLevelTo(logit.DebugLevel)
}
```

### 📖 参考案例

* [basic](./_examples/basic.go)
* [logger](./_examples/logger.go)
* [enable_disable](./_examples/enable_disable.go)
* [change_log_level](./_examples/change_log_level.go)

_更多使用案例请查看 [_examples](./_examples) 目录。_

### 🔥 性能测试

```bash
$ go test -v -bench=. -benchtime=20s
```

| 测试 | 单位时间内运行次数 (large is better) |  ns/op (small is better) | B/op (small is better) | allocs/op (small is better) |
| -----------|--------|-------------|-------------|-------------|
| **[logit](./logger_test.go)** | 4800000 | 5062 ns/op | 864 B/op | 8 allocs/op |
| [Golang log](./logger_test.go) | 5400000 | 4730 ns/op | 928 B/op | 12 allocs/op |

_由于目前的 logit 是基于 Golang log 的，所以成绩相比更差，后续会重新设计内部日志输出模块，所以当前成绩仅供参考！_

### 👥 贡献者

如果您觉得 logit 缺少您需要的功能，请不要犹豫，马上参与进来，发起一个 _**issue**_。

### 📦 使用 logit 的项目

| 项目 | 作者 | 描述 |
| -----------|--------|-------------|
|  |  |  |

