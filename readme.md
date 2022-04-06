<!--
 * @Author: zhangpanfeng
 * @Date: 2022-03-18 15:44:27
 * @LastEditors: zhangpanfeng
 * @LastEditTime: 2022-04-06 16:17:06
 * @FilePath: \get_feishu_uid\readme.md
 * email: 995199148@qq.com
 * Copyright (c) 2022 by zhangpanfeng, All Rights Reserved. 
-->
# 飞书open_id获取

* 主要功能：
  * 本程序主要是为了发送飞书webhook时可以@用户，获取用户open_id
* 项目原理：
  * 通过创建的应用，请求用户授权，将授权信息保存到redis，方便其他程序调用获取
* 使用前提：
  * 创建一个飞书自定义应用（普通用户即可创建，不上线也可以使用），获取到飞书的appid和appsecret保存到配置文件中
  * 在应用的安全设置中添加回调的安全域名接口（精确到接口），否则会被提示不安全未授权，域名可以通过更改hosts文件实现
  * 更改自己的redis和web地址配置


* 安装依赖：go mod tidy
* 运行方法：go run main.go
* 编译方法：go build
* 临时调试：air -c .air.conf
* 注意事项：配置文件config.yaml和编译后的二进制文件放同一级目录