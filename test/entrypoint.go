package main

import (
	"app/core"
	"fmt"
)

func main() {
	e := core.EntryPoint{
		Type:            core.TYPE_HTTP_PROXY,
		Name:            "demo",
		Domain:          "a.com",
		PathPrefixStrip: "/a",
		PassHostHeader:  true,
		LimitConn:       10,
		LimitRate:       "1m",
		HealthCheck:     &core.HealthCheck{5000, 10},
		Upstreams: []*core.Upstream{
			&core.Upstream{Url: "baidu.com", Weight: 10},
			&core.Upstream{Url: "google.com", Weight: 10},
		},
		WhiteList: []string{"127.0.0.1"},
		HttpAuthForward: &core.HttpAuthForward{
			Headers: []string{"HTTP-X-FORWARD"},
			Address: "b.com",
		},
		HttpRedirect: []*core.HttpRedirect{
			{Regex: "^http://a.com/(.*)", Replacement: "http://baidu.com/$1", Permanent: false},
		},
	}
	fmt.Println(e.Generate())
}
