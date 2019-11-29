package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"reflect"
	"time"
)

type ApiEtcd struct {
	addresses []string //etcd地址池
	autoClose bool     //是否自动关闭etcd连接

	_cli *clientv3.Client //连接实例
}

// 创建ApiEtcd实例
//options中bool类型代表是否为开启自动auto close etcd
//options中string类型代表是etcd address支持多个
func NewApiEtcd(options ...interface{}) *ApiEtcd {
	etcd := new(ApiEtcd)
	addresses := make([]string, 0)
	for _, v := range options {
		switch reflect.TypeOf(v).String() {
		case "string":
			addresses = append(addresses, v.(string))
		case "bool":
			etcd.autoClose = v.(bool)
		}
	}
	etcd.addresses = addresses
	return etcd
}

// 关闭etcd连接
func (this *ApiEtcd) Close() error {
	var err error
	func() {
		defer func() {
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprintf("%v", e))
			}
		}()
		if this._cli != nil {
			err = this._cli.Close()
		}
	}()
	return err
}

// 初始化etcd连接
func (this *ApiEtcd) cli() (*clientv3.Client, error) {
	if this._cli != nil {
		return this._cli, nil
	}

	if this.addresses == nil || len(this.addresses) == 0 {
		if len(this.addresses) == 0 {
			return nil, errors.New("etcd address is nil.")
		}
	}
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   this.addresses,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	this._cli = c
	return c, nil
}

// 获取某个key值
// params key: 键
// params isDir: 是否获取整个目录下的数据
func (this *ApiEtcd) GetKey(key string) ([]byte, error) {
	defer func() {
		if this.autoClose {
			this.Close()
		}
	}()

	c, err := this.cli()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var resp *clientv3.GetResponse
	resp, err = c.Get(ctx, key)
	cancel()
	if err != nil {
		return nil, err
	}
	for _, v := range resp.Kvs {
		return v.Value, nil
	}
	return []byte{}, nil
}

func (this *ApiEtcd) getKeysToBytes(key string) (map[string][]byte, error) {
	c, err := this.cli()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var resp *clientv3.GetResponse
	resp, err = c.Get(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, err
	}
	data := map[string][]byte{}
	for _, v := range resp.Kvs {
		data[string(v.Key)] = v.Value
	}
	return data, nil
}

// 获取某个key值下的所有内容
// params key: 键
func (this *ApiEtcd) GetKeys(key string) (map[string]string, error) {
	defer func() {
		if this.autoClose {
			this.Close()
		}
	}()

	res, err := this.getKeysToBytes(key)
	if err != nil {
		return nil, err
	}
	data := map[string]string{}
	for k, v := range res {
		data[k] = string(v)
	}
	return data, nil
}

// 修改某个key的值
// params key: 键
// params value: 值
// params options []interface{} 超时处理
func (this *ApiEtcd) UpdateKey(key, value string, options ...interface{}) error {
	defer func() {
		if this.autoClose {
			this.Close()
		}
	}()

	var ttl int64 = 0
	if len(options) > 0 {
		ttl = int64(options[0].(int))
	}
	c, err := this.cli()
	if err != nil {
		return err
	}
	var resp *clientv3.LeaseGrantResponse
	if ttl > 0 {
		resp, err = c.Grant(context.TODO(), ttl)
		if err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if resp != nil {
		_, err = c.Put(ctx, key, value, clientv3.WithLease(resp.ID))
	} else {
		_, err = c.Put(ctx, key, value)
	}
	cancel()
	return err
}

// 删除key
// params key: 键
// params isDir: 是否删除整个目录下的数据
func (this *ApiEtcd) DeleteKey(key string, isDir bool) error {
	defer func() {
		if this.autoClose {
			this.Close()
		}
	}()

	if key == "" {
		return errors.New("key is nil.")
	}
	c, err := this.cli()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if isDir {
		if key[len(key)-1:] != "/" {
			key += "/"
		}
		_, err = c.Delete(ctx, key, clientv3.WithPrefix())
	} else {
		_, err = c.Delete(ctx, key)
	}
	cancel()
	return err
}
