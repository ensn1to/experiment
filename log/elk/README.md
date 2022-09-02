# ELK

利用docker-compose方式部署elk


## 用法
### 1. 修改logstash/logstash.conf，配置需要的插件方式，已配置redis和tcp server方式


### 2. 给elasticsearch/data加写入权限
```
chmod -R 777 lasticsearch/data
```

### 3. 创建网络
此处创建bridge网络
```
docker create network  elk_default
```

如果单机开发调试，直接创建host网络更方便
```
service:
    images:redis
    ....
    network_mode: "host"
```

### 4. 启动
```
make upd
```

### 5. 关闭
```
make downd
```

> 更多方式参考Makefile
