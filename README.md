# CrestedIbis

## 功能列表

- GB28181:
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

## 依赖

- https://github.com/Monibuca 音视频处理
- https://github.com/gin-gonic/gin Web框架
- https://github.com/swaggo/gin-swagger Swagger API文档生成中间件
- https://github.com/spf13/cobra CLI应用程序框架
- https://github.com/sirupsen/logrus 日志框架
- https://github.com/casbin/casbin 权限管理框架
- https://github.com/golang-jwt/jwt JWT
- https://pkg.go.dev/golang.org/x/crypto 加盐加密
- https://github.com/go-yaml/yaml yaml文件
- https://github.com/go-gorm/gorm orm数据库处理

## 文件结构

# 文件结构

- dev: Docker开发依赖容器
- docs: Swagger API，`swag init`自动生成，无需受到编辑
- gb28181_server: Monibuca GB28181插件实现
- src: web服务实现
    - apps：app接口
    - config：定义config.yaml配置数据模型，初始化配置数据
    - global：全局变量、模型
- config.yaml：配置文件
