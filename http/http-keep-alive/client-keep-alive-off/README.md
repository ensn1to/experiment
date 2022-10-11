## 客户端关闭长连接

比如客户端在一条连接上请求密度并不高，不必长期占用链接资源


### 抓包
可见在client的GET请求头里多了Connection: close字段:
```
Frame 65: 191 bytes on wire (1528 bits), 191 bytes captured (1528 bits) on interface lo0, id 0
Null/Loopback
Internet Protocol Version 6, Src: ::1, Dst: ::1
Transmission Control Protocol, Src Port: 61076, Dst Port: 18081, Seq: 1, Ack: 1, Len: 115
Hypertext Transfer Protocol
    GET / HTTP/1.1\r\n
    Host: localhost:18081\r\n
    User-Agent: Go-http-client/1.1\r\n
    Accept-Encoding: gzip\r\n
    Connection: close\r\n
    \r\n
    [Full request URI: http://localhost:18081/]
    [HTTP request 1/1]
    [Response in frame: 67]

```

![抓包](./client_dis_keepalive.jpg)