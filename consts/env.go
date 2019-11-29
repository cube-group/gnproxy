package consts

const (
	ENV_PROVIDER         = "PROVIDER"             //配置源类型

	ENV_ETCD_ADDRESS     = "ETCD_ADDRESS"         //etcd provider address, 127.0.0.1:2379
	ENV_ETCD_HTTP_PREFIX = "ENV_ETCD_HTTP_PREFIX" //etcd provider prefix key name
	ENV_ETCD_TCP_PREFIX  = "ENV_ETCD_TCP_PREFIX"

	ENV_KUBE_API_SERVER = "KUBE_API_SERVER" //k8s api server
	ENV_KUBE_API_CONFIG = "KUBE_API_CONFIG" //k8s config
)
