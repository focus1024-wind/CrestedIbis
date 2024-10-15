package ipc

import (
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"CrestedIbis/src/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

// IpcUploadImage IPC图像上传
//
//	@Title			IPC图像上传
//	@Version		1.0.0
//	@Description	GB28181图像抓拍，图片上传接口
//	@Tags			IPC设备 /ipc/device
//	@Accept			mpfd
//	@Produce		json
//	@Param			access_token	query		string							true	"访问token"
//	@Param			file			formData	file							true	"上传图片"
//	@Success		200				{object}	model.HttpResponse{data=string}	"上传图片成功"
//	@Failure		200				{object}	model.HttpResponse{data=string}	"上传图片失败"
//	@Router			/ipc/device/upload_image [POST]
func IpcUploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		model.HttpResponse{}.FailGin(c, "获取图片数据失败")
		return
	}

	fileExt := strings.ToLower(path.Ext(file.Filename))
	if fileExt == ".jpg" || fileExt == ".png" || fileExt == ".jpeg" {
		filePath := global.Conf.Store.Snapshot
		claims, ok := c.Get("claims")
		if ok {
			switch jwtToken := claims.(type) {
			case utils.JwtToken:
				filePath = fmt.Sprintf("%s/%s", filePath, jwtToken.Username)
			}
		}

		err := c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", filePath, file.Filename))
		if err != nil {
			model.HttpResponse{}.FailGin(c, "保存图片失败")
			return
		} else {
			model.HttpResponse{}.OkGin(c, nil)
			return
		}
	} else {
		model.HttpResponse{}.FailGin(c, "图片格式错误，目前仅支持.jpg, .png, .jpeg格式数据上传")
		return
	}
}
