# ELK

利用docker-compose方式部署elk


## 用法
1. 修改logstash/logstash.conf，配置需要的插件方式，已配置redis和tcp server方式

2. 给elasticsearch/data加写入权限
```
chmod -R 777 lasticsearch/data
```

3. 启动
```
docker-compose up -d
```
