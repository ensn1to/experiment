## Docker logging-driver更改为fluentd

两种方式将容器日志加入到fluentd
1. 应用容器利用docker-compose.yml
```
services:
  web:
    image: httpd
    ports: 
      - 8080:80
    depends_on:
     - fluentd
    logging: *default-logging
```
2. docker run时加入
```
docker run -d \
  -p 8080:80 \
  --name httpd \
  --log-driver=fluentd \
  --log-opt fluentd-address=0.0.0.0:24224 \
  --log-opt mode=non-blocking \
  --log-opt tag={{.Name}} \
  --log-opt fluentd-async-connect \
  --network fluentd_default \
  httpd
```


ref:
1. https://www.notion.so/Docker-logging-driver-6463b3b5293f44ad8a1643a0edc591f7