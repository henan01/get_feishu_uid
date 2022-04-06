/*
 * @Author: 张攀锋
 * @Date: 2022-03-12 14:47:16
 * @LastEditors: 张攀锋
 * @LastEditTime: 2022-03-21 14:29:36
 * @FilePath: \get_feishu_uid\feishu\get_feishu_uid.go
 * email: 995199148@qq.com
 * Copyright (c) 2022 by zhangpanfeng, All Rights Reserved.
 */
package feishu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

const (
	// appid     = "cli_a131b044e879900c"
	// appsecret = "JrTKGF0UHMDeCW6ubXGAegeBDhqkeZxi"
	// // 重定向 URL，添加到，飞书开放平台安全设置》重定向URL中：https://open.feishu.cn/app/cli_a131b044e879900c/safe
	// redirect_uri = "http://getfeishuid.ktvsky.com:8888/uid"

	// 以下为飞书固定参数
	grant_type    = "authorization_code"
	get_code_url  = "https://passport.feishu.cn/suite/passport/oauth/authorize"
	get_token_url = "https://passport.feishu.cn/suite/passport/oauth/token"
	get_id_url    = "https://passport.feishu.cn/suite/passport/oauth/userinfo"
)

type AppConfig struct {
	Redis struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		Db   int    `mapstructure:"db"`
		Pwd  string `mapstructure:"pwd"`
	} `mapstructure:"redis"`
	Web struct {
		Host        string `mapstructure:"host"`
		RedirectURI string `mapstructure:"redirect_uri"`
	} `mapstructure:"web"`
	Feishu struct {
		Appid     string `mapstructure:"appid"`
		Appsecret string `mapstructure:"appsecret"`
	} `mapstructure:"feishu"`
}

var Conf = new(AppConfig)

func InitConfig() error {
	viper.SetConfigFile("config.yaml")

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("夭寿啦~配置文件被人修改啦...")
		viper.Unmarshal(&Conf)
	})

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(fmt.Errorf("unmarshal to Conf failed, err:%v", err))
	}
	return err
}

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func InitRedis() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Conf.Redis.Host, Conf.Redis.Port),
		Password: Conf.Redis.Pwd, // no password set
		DB:       Conf.Redis.Db,  // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// 批量向key的hash添加对应元素field的值
func BatchHashSet(key string, fields map[string]interface{}) string {
	val, err := rdb.HMSet(key, fields).Result()
	if err != nil {
		fmt.Println("Redis HMSet Error:", err)
	}
	return val
}

// 
func GetHashKey(key string, field string) (val string, err error) {
	val, err = rdb.HGet(key, field).Result()
	return
}

type TokenResult struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type UserInfoResult struct {
	Sub          string `json:"sub"`
	Picture      string `json:"picture"`
	Name         string `json:"name"`
	EnName       string `json:"en_name"`
	TenantKey    string `json:"tenant_key"`
	AvatarURL    string `json:"avatar_url"`
	AvatarThumb  string `json:"avatar_thumb"`
	AvatarMiddle string `json:"avatar_middle"`
	AvatarBig    string `json:"avatar_big"`
	OpenID       string `json:"open_id"`
	UnionID      string `json:"union_id"`
	Email        string `json:"email"`
}

func Get_code() (url string) {
	url = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code", get_code_url, Conf.Feishu.Appid, Conf.Web.RedirectURI)
	return
}

func Get_token(code string) TokenResult {
	// url = f'{get_token_url}?grant_type={grant_type}&client_id={appid}&client_secret={appsecret}&code={code}&redirect_uri={redirect_uri}'
	url := fmt.Sprintf("%s?grant_type=%s&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s", get_token_url, grant_type, Conf.Feishu.Appid, Conf.Feishu.Appsecret, code, Conf.Web.RedirectURI)
	res := TokenPostData(url)
	// fmt.Println(url)
	// fmt.Println(res)
	return res
}

func TokenPostData(url string) TokenResult {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Apipost client Runtime/+https://www.apipost.cn/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var tokenres TokenResult
	err = json.Unmarshal(bodyText, &tokenres)
	if err != nil {
		fmt.Println(err, "111111111111111111111")
	}
	return tokenres
}

func UserPostData(url, Token string) UserInfoResult {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Apipost client Runtime/+https://www.apipost.cn/")
	req.Header.Set("Authorization", "Bearer "+Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var userres UserInfoResult
	err = json.Unmarshal(bodyText, &userres)
	if err != nil {
		fmt.Println(err)
	}
	return userres
}

func Get_user_info(token string) (Userinfo UserInfoResult) {
	Userinfo = UserPostData(get_id_url, token)
	return
}
