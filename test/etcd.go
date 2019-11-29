package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main(){
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"security.proxy.eoffcn.com:23079"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	c.Put(context.Background(),"/gnproxy/backend/demo/passhostheader","true")
}
