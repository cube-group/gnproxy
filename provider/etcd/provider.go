package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type EtcdProvider struct {
}

func NewEtcdProvider() {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"security.proxy.eoffcn.com:23079"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	//var getResp *clientv3.GetResponse
	var watchResp clientv3.WatchResponse
	//var watchStartRevision int64
	var watchRespChan <-chan clientv3.WatchResponse
	var event *clientv3.Event
	key := "/gnproxy"
	// 先GET到当前的值，并监听后续变化
	//if _, err = c.Get(context.Background(), key); err != nil {
	//	fmt.Println(err)
	//	return
	//}
	// 当前etcd集群事务ID, 单调递增的（监听/cron/jobs/job7后续的变化,也就是通过监听版本变化）
	//watchStartRevision = getResp.Header.Revision + 1

	//watchRespChan = c.Watch(context.Background(), key, clientv3.WithRev(watchStartRevision))
	watchRespChan = c.Watch(context.Background(), key, clientv3.WithPrefix(), clientv3.WithPrevKV())
	// 处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			fmt.Printf("type:%v\n kv:%v  prevKey:%v  ", event.Type, event.Kv, event.PrevKv)
		}
	}

}
