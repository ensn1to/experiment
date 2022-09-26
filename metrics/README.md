## 应用和服务器监控

### How to use it
1. 在应用中配置要采样的数据，并启动一个http端口用于被在prometheus定时拉取采样

2. 在prometheus.yml配置对应的job
```
  - job_name: "tcp-server"
    static_configs:
      - targets: ["localhost:8889"]
```

3. 重启服务


注意：需要给data加写入权限:
```shell
chmod +x 777 ./data
```