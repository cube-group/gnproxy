package core

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"strings"
)

const (
	TYPE_HTTP_PROXY   = 0 //7层代理模式
	TYPE_STREAM_PROXY = 1 //4层代理模式
)

type EntryPoint struct {
	Name string //唯一标识

	Listen uint //监听端口

	LimitConn uint   //连接数限制
	LimitRate string //流量限制

	TimeoutConnect uint //连接超时时间，单位：秒，默认为4s
	TimeoutRead    uint //代理接收超时时间，单位：秒，默认为60s
	TimeoutSend    uint //代理发送超时时间，单位：秒，默认为12s

	HealthCheck *HealthCheck //tcp健康监测

	Upstreams []*Upstream //backend

	WhiteList []string //ip白名单
}

type HttpEntryPoint struct {
	IsWebSocket     bool
	IsGrpc          bool
	Host            string //接收的域或者ip
	PathPrefixStrip string //代理识别目录（代理需去除）
	PathPrefix      string //代理识别目录（无需去除）
	PassHostHeader  bool   //代理所有header

	HttpAuthForward *HttpAuthForward //auth forward
	HttpAuthBasic   HttpAuthBasic    //auth basic
	HttpRedirect    []*HttpRedirect  //跳转 http专有

	EntryPoint
}

type StreamEntryPoint struct {
	EntryPoint
}

type HealthCheck struct {
	FailTimeout uint //单位检测时间，单位秒，默认为15秒
	MaxFails    uint //单位检测时间内出现失败次数被认为不健康，默认为：2
}

type Upstream struct {
	Url    string
	Weight uint
}

type HttpAuthForward struct {
	Headers []string //forward后追加的header名称
	Address string   //forward地址
}

type HttpAuthBasic map[string]string

type HttpRedirect struct {
	Regex       string //location规则，例如：^http://localhost/(.*)
	Replacement string //rewrite，例如：http://mydomain/$1
	Permanent   bool   //是否追加permanent
}

//生成nginx配置文件

//生成nginx http反向代理配置文件
func (this *HttpEntryPoint) Generate() error {
	if this.Name == "" {
		return errors.New("Name is nil.")
	}

	//stream
	if this.Upstreams == nil {
		return errors.New("upstreams is nil.")
	}
	upstreamName := ""
	if this.IsWebSocket {
		upstreamName = "ws__" + this.Name
	} else if this.IsGrpc {
		upstreamName = "grpc__" + this.Name
	} else {
		upstreamName = "http__" + this.Name
	}
	upstream := "upstream " + upstreamName + "{\n"
	for _, s := range this.Upstreams {
		itemStream := "server " + s.Url + " "
		if this.HealthCheck != nil {
			itemStream += fmt.Sprintf(
				"max_fails=%d fail_timeout=%ds ",
				this.HealthCheck.MaxFails,
				this.HealthCheck.FailTimeout,
			)
		}
		if s.Weight > 0 {
			itemStream += fmt.Sprintf("weight=%d;\n", s.Weight)
		} else {
			itemStream += fmt.Sprintf(";\n", s.Weight)
		}
		upstream += itemStream;
	}
	upstream += "}\n"

	//server
	if this.Listen == 0 {
		this.Listen = 80
	}
	server := "server{\n"
	if this.IsGrpc {
		server += fmt.Sprintf("listen %d http2;\n", this.Listen)
	} else {
		server += fmt.Sprintf("listen %d;\n", this.Listen)
	}
	if this.Host == "" {
		return errors.New("domain is nil.")
	}
	server += "server_name " + this.Host + ";\n"

	if this.LimitConn > 0 {
		server += "limit_conn perserver_conn %d;\n"
		server = fmt.Sprintf(server, this.LimitConn)
	}
	if this.LimitRate != "" {
		server += "limit_rate %s;\n"
		server = fmt.Sprintf(server, this.LimitRate)
	}

	//代理核心
	var location string

	if this.HttpRedirect != nil {
		for _, r := range this.HttpRedirect {
			if r.Regex == "" || r.Replacement == "" {
				continue
			}
			route := r.Regex
			if strings.Index(r.Regex, "^http://"+this.Host) >= 0 { //跳转配置（兼容traefik v1)
				route = "/"
			}
			location += "location " + route + " {\n"
			location += "rewrite ^/(.*)$ " + r.Replacement
			if r.Permanent {
				location += " permanent;\n"
			} else {
				location += " redirect;\n"
			}
			location += "}\n"
		}
	}

	//核心路由
	proxyPass := ""
	if this.IsGrpc {
		proxyPass = "grpc_pass grpc://" + upstreamName + ";\n"
	} else if this.PathPrefix != "" {
		location += "location " + this.PathPrefix + " {\n"
		proxyPass = "proxy_pass http://" + upstreamName + ";\n"
	} else if this.PathPrefixStrip != "" {
		location += "location " + this.PathPrefixStrip + " {\n"
		proxyPass = "proxy_pass http://" + upstreamName + "/;\n"
	} else {
		location += "location / {\n"
		proxyPass = "proxy_pass http://" + upstreamName + "/;\n"
	}

	//ip whitelist
	if this.WhiteList != nil {
		for _, ip := range this.WhiteList {
			location += "allow " + ip + ";\n"
		}
		location += "deny all;\n"
	}

	//auth
	if this.HttpAuthForward != nil { //basic auth
		location += "auth_request /auth_request;\n"
		if this.HttpAuthForward.Headers != nil {
			for i, h := range this.HttpAuthForward.Headers {
				location += fmt.Sprintf("auth_request_set $forward_res_header_%d $upstream_http_%s;\n", i, h)
			}
			location += "\n"
			for i, h := range this.HttpAuthForward.Headers {
				location += fmt.Sprintf("proxy_set_header %s $forward_res_header_%d;\n", h, i)
			}
		}
	} else if this.HttpAuthBasic != nil { //forward auth
		basicAuth := ""
		for username, password := range this.HttpAuthBasic {
			if username == "" || password == "" {
				continue
			}
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				continue
			}
			basicAuth += string(hashedPassword) + "\n"
		}
		if err := ioutil.WriteFile("/etc/nginx/basic_auth/"+this.Name, []byte(basicAuth), 0777); err != nil {
			return err
		}
		location += "auth_basic gnproxy basic auth;\n";
		location += "auth_basic_user_file basic_auth/" + this.Name + ";\n";
	}

	//默认传递
	location += "proxy_set_header X-Forwarded-For $remote_addr;\n"
	if this.PassHostHeader {
		location += "proxy_set_header Host $host;\n"
	}
	location += "proxy_pass_request_headers on;\n"
	location += proxyPass
	//proxy_pass timeout setting
	if this.TimeoutConnect == 0 {
		this.TimeoutConnect = 4
	}
	location += fmt.Sprintf("proxy_connect_timeout %ds;\n", this.TimeoutConnect)
	if this.TimeoutRead == 0 {
		this.TimeoutRead = 60
	}
	location += fmt.Sprintf("proxy_read_timeout %ds;\n", this.TimeoutRead)
	if this.TimeoutSend == 0 {
		this.TimeoutSend = 4
	}
	location += fmt.Sprintf("proxy_send_timeout %ds;\n", this.TimeoutSend)
	//websocket支持
	if this.IsWebSocket {
		location += "proxy_http_version 1.1;\n"
		location += "proxy_set_header Upgrade $http_upgrade;\n"
		location += "proxy_set_header Connection \"Upgrade\";\n"
	}
	location += "}\n"

	//auth forward
	if this.HttpAuthForward != nil {
		location += "location /auth_request {\n"
		location += "internal;\n"
		location += "proxy_pass " + this.HttpAuthForward.Address + ";\n"
		location += "proxy_pass_request_body off;\n"
		location += "proxy_set_header X-Forwarded-For $remote_addr;\n"
		location += "proxy_set_header Host $host;\n"
		location += "proxy_set_header X-Original-URI $request_uri;\n"
		location += "}\n"
	}
	server += location
	server += "}\n"

	fmt.Println((upstream + server))

	cfgPath := "/etc/nginx/conf.d/" + this.Name + ".conf"
	if err := ioutil.WriteFile(
		cfgPath,
		[]byte(upstream+server),
		0777,
	); err != nil {
		return err
	}

	if err := Test(); err != nil {
		return err
	}

	if err := Reload(); err != nil {
		return err
	}

	return nil
}

//生成nginx tcp反向代理配置文件
func (this *StreamEntryPoint) Generate() error {
	if this.Name == "" {
		return errors.New("Name is nil.")
	}
	if this.Listen == 0 {
		this.Listen = 80
	}
	server := "server{\n"
	server += fmt.Sprintf("listen %d;\n", this.Listen)

	if this.LimitConn > 0 {
		server += "limit_conn ip_addr %d;\n"
		server = fmt.Sprintf(server, this.LimitConn)
	}
	if this.LimitRate != "" {
		server += "proxy_download_rate %s;\n"
		server += "proxy_upload_rate %s;\n"
		server = fmt.Sprintf(server, this.LimitRate, this.LimitRate)
	}

	if this.Upstreams == nil {
		return errors.New("upstreams is nil.")
	}
	upstream := "upstream stream__" + this.Name + "{\n"
	for _, s := range this.Upstreams {
		itemStream := "server " + s.Url + " "
		if this.HealthCheck != nil {
			itemStream += fmt.Sprintf(
				"max_fails=%d fail_timeout=%ds ",
				this.HealthCheck.MaxFails,
				this.HealthCheck.FailTimeout,
			)
		}
		if s.Weight > 0 {
			itemStream += fmt.Sprintf("weight=%d;\n", s.Weight)
		} else {
			itemStream += fmt.Sprintf(";\n", s.Weight)
		}
		upstream += itemStream;
	}
	upstream += "}\n"

	//ip whitelist
	if this.WhiteList != nil {
		for _, ip := range this.WhiteList {
			server += "allow " + ip + ";\n"
		}
		server += "deny all;\n"
	}
	server += "proxy_pass stream__" + this.Name + ";\n"
	server += "}\n"

	fmt.Println((upstream + server))

	cfgPath := "/etc/nginx/stream.d/" + this.Name + ".conf"
	if err := ioutil.WriteFile(
		cfgPath,
		[]byte(upstream+server),
		0777,
	); err != nil {
		return err
	}

	if err := Test(); err != nil {
		return err
	}

	if err := Reload(); err != nil {
		return err
	}

	return nil
}
