package initizalize

import (
	"CrestedIbis/docs"
	"CrestedIbis/src/apps/ipc"
	"CrestedIbis/src/apps/site"
	"CrestedIbis/src/apps/system"
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"CrestedIbis/src/utils"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"strings"
	"time"
)

// InitHttpServer 初始化HTTP服务器
func InitHttpServer() {
	gin.SetMode(gin.ReleaseMode)
	global.HttpEngine = gin.Default()

	corsConfig()

	global.HttpEngine.Use(ginLogger())
	global.HttpEngine.Use(permissionAuth())
	global.HttpEngine.Use(ginRecovery())

	docs.SwaggerInfo.BasePath = "/"
	// 设置路由组
	{
		ipc.InitIpcRouter()
		system.InitSystemRouter()
		site.InitSiteRouter()
	}
	global.HttpEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 优雅的关机: 监听关机信号，以释放系统资源
	// 详情参考: https://gin-gonic.com/zh-cn/docs/examples/graceful-restart-or-stop/
	// 在http_server仅启动Web程序，具体的事件监听和处理，由cobra进行处理
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", global.Conf.HttpServer.IP, global.Conf.HttpServer.Port),
		Handler: global.HttpEngine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.Logger.Fatalf("listen: %s", err)
		}
	}()
}

// cors跨域配置
func corsConfig() {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	global.HttpEngine.Use(cors.New(config))
}

// ginLogger 日志中间件配置
// 记录 HTTP 接口执行时间
func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// 设置HTTP日志信息
		global.Logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}

// ginRecovery 错误处理
func ginRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.Errorf("panic: %v", err)
				switch response := err.(type) {
				case model.HttpResponse:
					c.JSON(response.Code, response)
				case int:
					if response == http.StatusBadRequest {
						c.JSON(http.StatusBadRequest, model.HttpResponse{Code: http.StatusBadRequest, Msg: "error", Data: "参数错误"})
					}
				}
			}
			c.Abort()
		}()
		c.Next()
	}
}

// getIdentityRoles JWT 身份认证, 认证成功返回
// 认证方式1: Header: Authorization: Bear xxx(优先)
// 认证方式2: Params: access_token=xxx
// 默认权限：guest
func getIdentityRoles(c *gin.Context) ([]string, error) {
	tokenString := c.Request.Header.Get("Authorization")

	if len(strings.Split(tokenString, "Bearer ")) == 2 {
		tokenString = strings.Split(tokenString, "Bearer ")[1]
	} else if tokenString == "" {
		tokenString = c.DefaultQuery("access_token", "")
	}

	if tokenString != "" {
		claims, err := utils.JwtToken{}.ParseToken(tokenString)
		if err != nil {
			return []string{"guest"}, err
		}
		c.Set("claims", claims)
		return append(claims.Roles, "guest"), nil
	} else {
		return []string{"guest"}, nil
	}
}

// permissionAuth 通过casbin进行权限认证
func permissionAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 获取用户权限列表, 错误信息后续处理，避免过期Token重新调用登录接口
		roles, roles_err := getIdentityRoles(c)

		// 获取casbin鉴权执行器
		casbinEnforcer := utils.CasbinService()
		if casbinEnforcer == nil {
			// 初始化权限执行器出错
			c.JSON(http.StatusInternalServerError, model.HttpResponse{
				Code: http.StatusInternalServerError,
				Msg:  "error",
				Data: "权限校验失败",
			})
			c.Abort()
			return
		}

		// 对权限列表进行校验
		permissionRes := false
		for _, role := range roles {
			success, err := casbinEnforcer.Enforce(role, c.Request.URL.Path, c.Request.Method)
			if err != nil {
				// 检测权限出错
				c.JSON(http.StatusInternalServerError, model.HttpResponse{
					Code: http.StatusInternalServerError,
					Msg:  "error",
					Data: err.Error(),
				})
				c.Abort()
				return
			}
			if success {
				permissionRes = true
				c.Next()
				break
			}
		}

		if roles_err != nil && !permissionRes {
			if errors.Is(roles_err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, model.HttpResponse{
					Code: http.StatusUnauthorized,
					Msg:  "error",
					Data: "登录已过期",
				})
			} else {
				c.JSON(http.StatusUnauthorized, model.HttpResponse{
					Code: http.StatusUnauthorized,
					Msg:  "error",
					Data: roles_err.Error(),
				})
			}
			c.Abort()
			return
		}

		if !permissionRes {
			c.JSON(http.StatusForbidden, model.HttpResponse{
				Code: http.StatusForbidden,
				Msg:  "error",
				Data: "无访问权限",
			})
			c.Abort()
			return
		}
	}
}
