# gnproxy
Simple and fast cloud native network agent

<p>
  <img alt="gnproxy" height=100 width=150 src="img/gnproxy.png">
</p>

<p>
  <a href="https://search.maven.org/artifact/com.easymybatis.freamwork/spring-easymybatis-core">
    <img alt="maven" src="https://img.shields.io/badge/golang-1.11-blue">
  </a>
  <a href="https://github.com/996icu/996.ICU/blob/master/LICENSE">
    <img alt="996icu" src="https://img.shields.io/badge/license-NPL%20(The%20996%20Prohibited%20License)-blue.svg">
  </a>

  <a href="https://github.com/onlyGuo/easymybatis/blob/master/LICENSE">
    <img alt="code style" src="https://img.shields.io/badge/license-Apache%202-blue">
  </a>
</p>

Using golang to develop lightweight network agent based on nginx

### Features
* Our upper layer uses golang for development
* The bottom layer depends on nginx core
* Support TCP, websocket, grpc and HTTP protocol proxy
* Supports frontend and backend normal agents
* Support kubernetes label mode to configure agent
* Supports etcd kV configuration mode and is compatible with <a href='https://github.com/containous/traefik' target='_blank'>traefik</a> V1 etcd
* Support current limiting, fusing and simple statistics
* More features please look forward to :)