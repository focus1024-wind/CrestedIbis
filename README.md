# CrestedIbis

CrestedIbis目前是一个基于GB28181标准实现的音视频云平台，负责实现GB28181信令和设备管理，未来将会是一个支持物联网设备接入，算法训练和部署的综合物联网平台。

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<br />

<div style="text-align: center;">
<h3 style="text-align: center;">CrestedIbis</h3>
  <p style="text-align: center;">
    GB28181音视频平台!
    <br />
    <a href="https://github.com/focus1024-wind/CrestedIbis"><strong>探索本项目的文档 »</strong></a>
    <br />
    <br />
    <a href="https://github.com/focus1024-wind/CrestedIbis">查看Demo</a>
    ·
    <a href="https://github.com/focus1024-wind/CrestedIbis/issues">报告Bug</a>
    ·
    <a href="https://github.com/focus1024-wind/CrestedIbis/issues">提出新特性</a>
  </p>
</div>

## 目录

- [上手指南](#上手指南)
    - [开发前的配置要求](#开发前的配置要求)
    - [安装步骤](#安装部署)
- [文件目录说明](#文件目录说明)
- [使用到的框架](#使用到的框架)
- [gb28181功能列表](#gb28181功能列表)
- [贡献者](#贡献者)
- [license](#license)

### 上手指南

#### 开发前的配置要求

1. golang 1.23.0(
   作者当前开发环境，对于其他golang版本不做兼容性测试，若有不支持版本，请提[ISSUES](https://github.com/focus1024-wind/CrestedIbis/issues))
2. 数据库(根据你的喜好选择你的数据库类型)：
    - MySQL
    - PostgreSQL

#### 安装部署

1. 准备数据库环境，如果你有Docker环境，请参考如下步骤。如果没有，请搜索相关文档，自己安装数据库，本文档不做介绍
    ```shell
    cd dev
    docker compose up -d
    ```
2. 克隆本仓库
    ```shell
    git clone https://github.com/focus1024-wind/CrestedIbis.git
    ```
3. 下载依赖
    ```shell
    go mod download
    ```
4. 启动服务
    ```shell
    # 读取当前路径下config.yaml文件，运行
    go run main.go
    # 读取指定config文件
    go run main.go -c ${config_file}
    ```
5. 准备数据库数据，执行resources数据库脚本(本项目依赖[gorm](https://github.com/go-gorm/gorm)，启动时，会自动创建数据库表)

> [本项目依赖Casbin](https://github.com/casbin/casbin)进行权限管理，所以在首次初始化数据库表后，一定要执行数据库脚本，否则相关接口将无法访问。

### 文件目录说明

```
filetree
├─dev                   # 开发环境，Docker开发compose脚本
├─docs                  # Swagger API文档，`swag init`自动生成，无需手动编辑
├─gb28181_server        # 基于Monibuca实现的GB28181插件
├─resources             # 资源文档，存放数据库脚本
├─src                   # web服务实现
│  ├─apps               # app接口
│  │  └─ipc
│  ├─config             # 读取配置文件数据
│  │  └─model
│  ├─global             # 全局配置变量，模型
│  │  ├─initizalize     # 初始化脚本
│  │  └─model           # 全局模型
│  └─utils              # 工具类
├─store                 # 默认资源存储路径(图片抓拍，视频录制)
├─config.yaml           # 默认配置文件
└─main.go               # 主程序
```

### 使用到的框架

- [Monibuca音视频处理](https://github.com/Monibuca)
- [gin Web框架](https://github.com/gin-gonic/gin)
- [Swagger API文档生成](https://github.com/swaggo/gin-swagger)
- [cobra CLI应用程序框架](https://github.com/spf13/cobra)
- [logrusr日志框架](https://github.com/sirupsen/logrus)
- [casbin权限管理框架](https://github.com/casbin/casbin)
- [gorm orm数据库处理](https://github.com/go-gorm/gorm)

### GB28181功能列表

- GB28181-2022:
    - [ ] 9.1 注册和注销
        - [x] 9.1.2.1 基本注册
        - [x] 9.1.2.2 基本注销
        - [ ] 9.1.2.3 注册重定向
    - [ ] 9.2 实时音视频点播
        - [x] 9.2.2.1 客户端主动发起
        - [ ] 9.2.2.2 第三方呼叫控制
    - [ ] 9.3 控制
        - [ ] 9.3.2.1 无应答命令流程
        - [ ] 9.3.2.2 有应答命令流程
    - [ ] 9.4 报警事件通知和分发
    - [ ] 9.5 网络设备信息查询
    - [ ] 9.6 状态信息报送
    - [ ] 9.7 设备视音频文件检索
    - [ ] 9.8 历史视音频回放
        - [ ] 9.8.2.1 客户端主动发起
        - [ ] 9.8.2.2 第三方呼叫控制
    - [ ] 9.9 视音频文件下载
        - [ ] 9.9.2.1 客户端主动发起
        - [ ] 9.9.2.2 第三方呼叫控制
    - [x] 9.10 校时
    - [ ] 9.11 订阅和通知
        - [ ] 9.11.1.2 事件订阅
        - [ ] 9.11.2.2 事件通知
        - [ ] 9.11.3.2 目录订阅
        - [ ] 9.11.4.2 目录通知
    - [ ] 9.12 语音广播和语音对讲
        - [ ] 9.12.1.2 语音广播
        - [ ] 9.12.2.2 语音对讲
    - [ ] 9.13 设备软件升级
    - [x] 9.14 图像抓拍

## 贡献者

<a href="https://github.com/focus1024-wind/CrestedIbis/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=focus1024-wind/CrestedIbis"  alt="contributors"/>
</a>

## License

Apache License, Version 2.0, ([LICENSE-APACHE](LICENSE-APACHE)
or [http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)).