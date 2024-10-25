package gb28181_server

import (
	"CrestedIbis/gb28181_server/utils"
	"crypto/md5"
	"fmt"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const RegisterTimeLayout = time.DateTime

var (
	// DeviceNonce 存储设备注册 Nonce 信息、防止伪造
	DeviceNonce sync.Map
	// DeviceKeepalive 存储设备注册信息，key：deviceID，value：RegisterTime
	DeviceKeepalive sync.Map
	// DeviceChannels 存储设备通道ID，避免自动拉流，多次搜索数据库，减少数据库压力
	DeviceChannels sync.Map
)

func (config *GB28181Config) SipRegisterHandler(req sip.Request, tx sip.ServerTransaction) {
	if deviceID, ok := getGB28181DeviceIdBySip(req); ok {
		if expiresValue, ok := handlerExpiresValue(req); ok {
			if expiresValue != 0 {
				// 注册设备
				if config.Password != "" {
					// 有密码参数，多次认证登陆
					if authHeaders := req.GetHeaders("Authorization"); len(authHeaders) > 0 {
						// 有Authorization请求头，第二次上报
						// 第二次上报，携带Authorization信息，进行密码校验
						authHeader := authHeaders[0].(*sip.GenericHeader)
						auth := &Authorization{Authorization: sip.AuthFromValue(authHeader.Contents)}

						var username string
						if auth.Username() == deviceID {
							username = deviceID
						} else {
							username = config.Username
						}

						nonce, ok := DeviceNonce.Load(deviceID)
						if ok && auth.Verify(username, config.Password, config.Realm, nonce.(string)) {
							// 校验成功，注册设备
							registerDevice(deviceID, req, tx)
						} else {
							// 校验失败
							logger.Errorf("[SIP SERVER] DeviceID: %s 设备认证失败", deviceID)
							return
						}
					} else {
						// 无Authorization请求头，第一次上报
						// 第一次上报，返回WWW-Authorization
						unAuthorization(deviceID, req, tx)
					}
				} else {
					// 不需要校验：直接注册设备
					registerDevice(deviceID, req, tx)
				}
			} else {
				// 注销设备
				// Expires为0：注销设备
				logoutDevice(deviceID, req, tx)
				return
			}
		} else {
		}
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, "BAD REQUEST", ""))
	} else {
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusBadRequest, "BAD REQUEST", ""))
		return
	}
}

// handlerExpiresValue 获取REGISTER Expires头信息
// >0: 注册设备
// =0: 注销设备
// 其他: 异常
func handlerExpiresValue(req sip.Request) (expireValue int64, ok bool) {
	if expiresHeader := req.GetHeaders("Expires"); len(expiresHeader) > 0 {
		// 获取expires头，并进行格式化
		expireValue, err := strconv.ParseInt(expiresHeader[0].Value(), 10, 32)
		if err == nil && expireValue >= 0 {
			return expireValue, true
		}
	}
	return -1, false
}

// logoutDevice 注销设备
func logoutDevice(deviceId string, req sip.Request, tx sip.ServerTransaction) {
	GlobalGB28181DeviceStore.DeviceOffline(deviceId)

	response := sip.NewResponseFromRequest("", req, http.StatusOK, "OK", "")

	to, _ := response.To()
	response.ReplaceHeaders("To", []sip.Header{
		&sip.ToHeader{
			Address: to.Address,
			Params: sip.NewParams().Add(
				"tag", sip.String{
					Str: utils.RandNumString(9),
				},
			),
		},
	})

	response.RemoveHeader("Allow")

	expires := sip.Expires(3600)
	response.AppendHeader(&expires)

	response.AppendHeader(&sip.GenericHeader{
		HeaderName: "Date",
		Contents:   time.Now().Format(RegisterTimeLayout),
	})

	_ = tx.Respond(response)
}

// unAuthorization
// 设备未认证：返回 401状态码，WWW-Authenticate 消息头
func unAuthorization(deviceId string, req sip.Request, tx sip.ServerTransaction) {
	// 返回WWW-Authorization
	logger.Infof("[SIP SERVER] DeviceID: %s 设备未认证", deviceId)

	response := sip.NewResponseFromRequest("", req, http.StatusUnauthorized, "StatusUnauthorized", "")
	nonce, _ := DeviceNonce.LoadOrStore(deviceId, utils.RandNumString(32))
	auth := fmt.Sprintf(`Digest realm="%s",algorithm=%s,nonce="%s"`, globalGB28181Config.Realm, "MD5", nonce.(string))
	response.AppendHeader(&sip.GenericHeader{
		HeaderName: "WWW-Authenticate",
		Contents:   auth,
	})
	_ = tx.Respond(response)
}

// registerDevice 认证成功，设备注册
func registerDevice(deviceId string, req sip.Request, tx sip.ServerTransaction) {
	// 存储设备信息
	var device GB28181Device
	device.storeDevice(req, true)

	DeviceNonce.Delete(deviceId)
	DeviceKeepalive.Store(deviceId, time.Now())
	logger.Infof("[SIP SERVER] DeviceID: %s 设备注册", deviceId)

	// 注册响应
	response := sip.NewResponseFromRequest("", req, http.StatusOK, "OK", "")

	to, _ := response.To()
	response.ReplaceHeaders("To", []sip.Header{
		&sip.ToHeader{
			Address: to.Address,
			Params: sip.NewParams().Add(
				"tag", sip.String{
					Str: utils.RandNumString(9),
				},
			),
		},
	})

	response.RemoveHeader("Allow")

	expires := sip.Expires(3600)
	response.AppendHeader(&expires)

	response.AppendHeader(&sip.GenericHeader{
		HeaderName: "Date",
		Contents:   time.Now().Format(RegisterTimeLayout),
	})

	_ = tx.Respond(response)

	// 同步通道目录信息
	go device.syncChannels()
}

type Authorization struct {
	*sip.Authorization
}

// Verify 验证请求头 Authorization
func (authorization *Authorization) Verify(username, passwd, realm, nonce string) bool {
	// HA1 = MD5(username:realm:password)
	hash1 := authorization.encryption(fmt.Sprintf("%s:%s:%s", username, realm, passwd))
	// HA2 = MD5(method:digestURI)
	hash2 := authorization.encryption(fmt.Sprintf("REGISTER:%s", authorization.Uri()))

	var cipherText string
	if authorization.Qop() == "" {
		// cipherText = MD5(HA1:nonce:HA2)
		plainText := fmt.Sprintf("%s:%s:%s", hash1, nonce, hash2)
		cipherText = authorization.encryption(plainText)
	} else if authorization.Qop() == "auth" || authorization.Qop() == "auth-int" {
		// cipherText = MD5(HA1:nonce:nc:cnonce:qop:HA2)
		plainText := fmt.Sprintf("%s:%s:%s:%s:%s:%s", hash1, nonce, authorization.Nc(), authorization.CNonce(), authorization.Qop(), hash2)
		cipherText = authorization.encryption(plainText)
	} else {
		logger.Errorf("非法登陆鉴权算法")
		return false
	}

	return cipherText == authorization.Response()
}

// 加密
func (authorization *Authorization) encryption(raw string) string {
	switch authorization.Algorithm() {
	case "MD5":
		return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
	default:
		//如果没有算法，默认使用MD5
		return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
	}
}
