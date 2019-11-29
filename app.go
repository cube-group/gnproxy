package main

import (
	"app/conf"
	"app/core"
	"fmt"
	"time"
)

func main() {
	fmt.Println("hello gnproxy :)")

	conf.Init()

	//http
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("--->http")
		e := core.HttpEntryPoint{
			Host:            "a.com",
			PathPrefixStrip: "/a",
			PassHostHeader:  true,
			HttpAuthForward: &core.HttpAuthForward{
				Headers: []string{"HTTP-X-FORWARD"},
				Address: "http://b.com",
			},
			HttpRedirect: []*core.HttpRedirect{
				{Regex: "^http://a.com/(.*)", Replacement: "http://baidu.com/$1", Permanent: false},
			},
		}
		e.Name = "demo"
		e.LimitConn = 10
		e.LimitRate = "1m"
		e.HealthCheck = &core.HealthCheck{5000, 10}
		e.Upstreams = []*core.Upstream{
			{Url: "baidu.com", Weight: 10},
			{Url: "google.com", Weight: 10},
		}
		e.WhiteList = []string{"127.0.0.1"}
		fmt.Println(e.Generate())
	}()

	//tcp
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("--->tcp")
		e := core.StreamEntryPoint{

		}
		e.Listen = 10000
		e.Name = "demo2"
		e.LimitConn = 10
		e.LimitRate = "1m"
		e.HealthCheck = &core.HealthCheck{5000, 10}
		e.Upstreams = []*core.Upstream{
			{Url: "baidu.com:80", Weight: 10},
			{Url: "google.com:80", Weight: 10},
		}
		e.WhiteList = []string{"127.0.0.1"}
		//fmt.Println(e.Generate())
	}()

	//websocket
	go func() {
		time.Sleep(8 * time.Second)
		fmt.Println("--->ws")
		e := core.HttpEntryPoint{
			IsWebSocket: true,
			Host:            "a.com",
			PathPrefixStrip: "/a",
			PassHostHeader:  true,
			HttpAuthForward: &core.HttpAuthForward{
				Headers: []string{"HTTP-X-FORWARD"},
				Address: "http://b.com",
			},
			HttpRedirect: []*core.HttpRedirect{
				{Regex: "^http://a.com/(.*)", Replacement: "http://baidu.com/$1", Permanent: false},
			},
		}
		e.Listen = 10000
		e.Name = "demo3"
		e.LimitConn = 10
		e.LimitRate = "1m"
		e.HealthCheck = &core.HealthCheck{50, 10}
		e.Upstreams = []*core.Upstream{
			{Url: "baidu.com:80", Weight: 10},
		}
		e.WhiteList = []string{"127.0.0.1"}
		//fmt.Println(e.Generate())
	}()

	//grpc
	go func() {
		time.Sleep(11 * time.Second)
		fmt.Println("--->grpc")
		e := core.HttpEntryPoint{
			IsGrpc: true,
			Host:            "a.com",
			PathPrefixStrip: "/a",
			PassHostHeader:  true,
			HttpAuthForward: &core.HttpAuthForward{
				Headers: []string{"HTTP-X-FORWARD"},
				Address: "http://b.com",
			},
			HttpRedirect: []*core.HttpRedirect{
				{Regex: "^http://a.com/(.*)", Replacement: "http://baidu.com/$1", Permanent: false},
			},
		}
		e.Listen = 10000
		e.Name = "demo4"
		e.LimitConn = 10
		e.LimitRate = "1m"
		e.HealthCheck = &core.HealthCheck{50, 10}
		e.Upstreams = []*core.Upstream{
			{Url: "baidu.com:80", Weight: 10},
		}
		e.WhiteList = []string{"127.0.0.1"}
		//fmt.Println(e.Generate())
	}()

	core.Exec(func(row string) {
		fmt.Println(row)
	})

	time.Sleep(time.Hour)

}
