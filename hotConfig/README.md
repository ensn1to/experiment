## 热加载

利用viper包动态监控配置文件变更.
<br>
</br>

启动服务，修改配置：
```
➜  hotConfig git:(master) ✗ go run main.go
prepare host : host:127.0.0.1 port:6379

检测到配置更改...
filechange host : host:127.0.0.1 port:637999
检测到配置更改...
filechange host : host:127.0.0.1 port:6379
```