/*
 * @Author: 张攀锋
 * @Date: 2022-03-11 21:06:05
 * @LastEditors: 张攀锋
 * @LastEditTime: 2022-03-21 15:16:42
 * @FilePath: \get_feishu_uid\main.go
 * email: 995199148@qq.com
 * Copyright (c) 2022 by zhangpanfeng, All Rights Reserved.
 */
package main

import (
	"encoding/json"
	"fmt"
	"get_feishu_uid/feishu"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, feishu.Get_code())
	})

	r.GET("/uid", func(c *gin.Context) {
		code := c.Query("code")
		tk := feishu.Get_token(code)
		// 判断code是否被使用过，使用过，重定向到 /
		if tk.ErrorDescription == "The code is invalid." {
			c.Redirect(302, feishu.Get_code())
		}
		user := feishu.Get_user_info(tk.AccessToken)

		data, _ := json.Marshal(&user)
		m := make(map[string]interface{})
		json.Unmarshal(data, &m)

		feishu.BatchHashSet(user.Email, m)
		c.JSON(200, gin.H{
			"openid": user.OpenID,
		})
	})

	r.GET("/user", func(c *gin.Context) {
		// 通过邮箱地址，返回at 的open_id
		emailsStr := c.Query("emails")
		emailsList := strings.Split(emailsStr, ",")
		data := ""
		for _, v := range emailsList {
			open_id, err := feishu.GetHashKey(v, "open_id")
			if err != nil {
				fmt.Println(err)
			}
			if open_id != "" {
				data += fmt.Sprintf("<at user_id='%s'></at>", open_id)
			}

		}
		c.PureJSON(200, data)
	})

	feishu.InitConfig()
	feishu.InitRedis()
	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run(feishu.Conf.Web.Host)
}
