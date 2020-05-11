// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/agent/bind": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "代理"
                ],
                "summary": "绑定上级",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "agent_id",
                        "name": "agent_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/bank": {
            "get": {
                "consumes": [
                    "application/json; charset=utf-8"
                ],
                "tags": [
                    "银行卡"
                ],
                "summary": "获取银行卡信息",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/bink/bind": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "银行卡"
                ],
                "summary": "绑定银行卡",
                "parameters": [
                    {
                        "description": "请求的json结构",
                        "name": "json",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.RequestBankCardInfo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/check": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "审核"
                ],
                "summary": "用户获取提交的审核信息",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "审核"
                ],
                "summary": "代练提交或更新段位审核",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "game_id",
                        "name": "game_id",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "game_server",
                        "name": "game_server",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "game_pos",
                        "name": "game_pos",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "game_level",
                        "name": "game_level",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "img_url",
                        "name": "img_url",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/check/admin": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "审核"
                ],
                "summary": "Get 管理员获取审核列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/check/admin/{user_id}": {
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "审核"
                ],
                "summary": "管理员进行审核",
                "parameters": [
                    {
                        "description": "State -1/1",
                        "name": "state",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/order": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "下单"
                ],
                "summary": "下单",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "游戏",
                        "name": "game_type",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "订单类型",
                        "name": "order_type",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "代练类型",
                        "name": "instead_type",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "游戏区服",
                        "name": "game_zone",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "铭文等级",
                        "name": "runes_level",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "英雄数量",
                        "name": "hero_num",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "当前段位",
                        "name": "cur_level",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "目标段位",
                        "name": "target_level",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "游戏账号",
                        "name": "game_account",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "游戏密码",
                        "name": "game_pwd",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "游戏角色名",
                        "name": "game_role",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "验证手机",
                        "name": "game_phone",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "保证金",
                        "name": "margin",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "有防沉迷",
                        "name": "anti_addiction",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "有指定英雄",
                        "name": "designate_hero",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "指定英雄",
                        "name": "hero_name",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/order/all": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "接单"
                ],
                "summary": "Get 获取订单列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/order/confirm": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "订单"
                ],
                "summary": "完成订单",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "order_id",
                        "name": "order_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/order/finish": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "订单"
                ],
                "summary": "完成订单",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "order_id",
                        "name": "order_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/order/take": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "接单"
                ],
                "summary": "接单",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "order_id",
                        "name": "order_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/pay": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "微信支付"
                ],
                "summary": "微信下单 获取发起微信支付所需的数据",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "order_id",
                        "name": "order_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/pay/taker": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "微信支付"
                ],
                "summary": "微信接单 获取发起微信支付所需的数据",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "order_id",
                        "name": "order_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/phone/bind": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "绑定手机号"
                ],
                "summary": "绑定手机号 确认绑定的接口",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "phone",
                        "name": "phone",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "code",
                        "name": "code",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/phone/code": {
            "get": {
                "consumes": [
                    "application/json; charset=utf-8"
                ],
                "tags": [
                    "绑定手机号"
                ],
                "summary": "获取要绑定的手机号的验证码 绑定手机时使用",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "phone",
                        "name": "phone",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/phone/code2": {
            "get": {
                "consumes": [
                    "application/json; charset=utf-8"
                ],
                "tags": [
                    "手机验证码"
                ],
                "summary": "获取已绑定手机号的验证码 （不需要传手机号 验证提款密码时使用）",
                "parameters": [
                    {
                        "description": "user_id",
                        "name": "user_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "tags": [
                    "登陆"
                ],
                "summary": "微信登陆接口 发送code获取session_key",
                "parameters": [
                    {
                        "description": "sessionKey",
                        "name": "session_key",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/upload": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "审核"
                ],
                "summary": "上传段位审核截图",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Image File",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/wxlogin": {
            "post": {
                "tags": [
                    "登陆"
                ],
                "summary": "微信登陆接口 发送code获取session_key",
                "parameters": [
                    {
                        "description": "code",
                        "name": "code",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "nickname",
                        "name": "nickname",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "avatar_url",
                        "name": "avatar_url",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "gender",
                        "name": "gender",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "province",
                        "name": "province",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "city",
                        "name": "city",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "models.BankCardInfo": {
            "type": "object",
            "properties": {
                "bankBranchName": {
                    "description": "支行名称",
                    "type": "string"
                },
                "bankCard": {
                    "description": "银行卡号",
                    "type": "string"
                },
                "bankName": {
                    "description": "银行名称",
                    "type": "string"
                },
                "password": {
                    "description": "给银行卡信息设置的密码",
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                },
                "userName": {
                    "description": "开户名",
                    "type": "string"
                }
            }
        },
        "models.RequestBankCardInfo": {
            "type": "object",
            "properties": {
                "bankCardInfo": {
                    "type": "object",
                    "$ref": "#/definitions/models.BankCardInfo"
                },
                "code": {
                    "description": "验证码",
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Golang Gin API",
	Description: "A web of gin",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
