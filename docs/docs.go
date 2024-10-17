// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "CrestedIbis",
            "url": "https://github.com/focus1024-wind/CrestedIbis",
            "email": "focus1024@foxmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ipc/device/upload_image": {
            "post": {
                "description": "GB28181图像抓拍，图片上传接口",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "IPC设备 /ipc/device"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "访问token",
                        "name": "access_token",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "上传图片",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "上传图片失败",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.HttpResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/system/user/login": {
            "post": {
                "description": "用户登录并生成用户登录日志信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理 /system/user"
                ],
                "parameters": [
                    {
                        "description": "用户登录信息，密码采用加盐加密",
                        "name": "SysLoginUser",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.SysUserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "登录成功，响应JWT",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.HttpResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "登录失败，响应失败信息",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.HttpResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/system/user/register": {
            "post": {
                "description": "注册用户",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户管理 /system/user"
                ],
                "parameters": [
                    {
                        "description": "用户注册信息，密码采用加盐加密",
                        "name": "SysUser",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.SysUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "注册成功",
                        "schema": {
                            "$ref": "#/definitions/model.HttpResponse"
                        }
                    },
                    "500": {
                        "description": "注册失败，响应失败信息",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.HttpResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "gorm.DeletedAt": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "valid": {
                    "description": "Valid is true if Time is not NULL",
                    "type": "boolean"
                }
            }
        },
        "model.HttpResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "msg": {
                    "type": "string"
                }
            }
        },
        "user.SysUser": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "created_time": {
                    "type": "string"
                },
                "deletedAt": {
                    "$ref": "#/definitions/gorm.DeletedAt"
                },
                "email": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "example": "CrestedIbis"
                },
                "phone": {
                    "type": "string"
                },
                "sex": {
                    "type": "integer"
                },
                "updated_time": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                },
                "username": {
                    "type": "string",
                    "example": "admin"
                }
            }
        },
        "user.SysUserLogin": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "example": "CrestedIbis"
                },
                "username": {
                    "type": "string",
                    "example": "admin"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
