## 客户端管理多个服务端连接

一个http client可管理到多个server的连接，并优先复用到同一server的连接(keep-alive)，而不是建立新连


### 流程
启两个http server，一个http client， client同时发请求给两个server，并连续发两次，观察连续发两次是否是服用连接的

### 抓包
![抓包](./client-mulit-server.png)

