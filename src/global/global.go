package global

import (
	"CrestedIbis/src/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 全局变量
var (
	Conf       *config.Config
	Logger     *logrus.Logger
	Db         *gorm.DB
	HttpEngine *gin.Engine
)
