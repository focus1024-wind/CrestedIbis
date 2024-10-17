package initizalize

import (
	"CrestedIbis/docs"
	"CrestedIbis/src/apps/ipc"
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
	}
	global.HttpEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err := global.HttpEngine.Run(fmt.Sprintf("%s:%d", global.Conf.HttpServer.IP, global.Conf.HttpServer.Port))
	if err != nil {
		global.Logger.Printf("http server start error: %s", err)
	}
}

// cors跨域配置
func corsConfig() {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	global.HttpEngine.Use(cors.New(config))
}

// ginLogger 日志中间件配置
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
				switch response := err.(type) {
				case model.HttpResponse:
					c.JSON(response.Code, response)
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

		// 获取用户权限列表
		roles, err := getIdentityRoles(c)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, model.HttpResponse{
					Code: http.StatusUnauthorized,
					Msg:  "error",
					Data: "登录已过期",
				})
			} else {
				c.JSON(http.StatusUnauthorized, model.HttpResponse{
					Code: http.StatusUnauthorized,
					Msg:  "error",
					Data: err.Error(),
				})
			}
			c.Abort()
			return
		}

		// 获取casbin鉴权执行器
		casbinEnforcer := utils.CasbinService()
		if casbinEnforcer == nil {
			// 初始化权限执行器出错
			c.JSON(http.StatusInternalServerError, model.HttpResponse{
				Code: http.StatusInternalServerError,
				Msg:  "权限校验失败",
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
					Msg:  err.Error(),
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
		if !permissionRes {
			c.JSON(http.StatusForbidden, model.HttpResponse{
				Code: http.StatusForbidden,
				Msg:  "无访问权限",
			})
			c.Abort()
			return
		}
	}
}
