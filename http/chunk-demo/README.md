## http chunk request demo

### 服务端分块规则
- 每个分块包含两个部分，长度头和数据块；
- 长度头是以CRLF（回车换行，即\r\n）结尾的一行明文，用16进制数字表示长度；
- 数据块紧跟在长度头后，最后也用CRLF结尾，但数据不包含CRLF；
- 最后用一个长度为0的块表示结束，即“0\r\n\r\n”

more: https://www.notion.so/Cornell-Notes-System-4b9b4e8b0529401da34b4ad6375f9f65#741890cdc74b49f393c963add1c8121f

<br>
</br>

### 实现
<p>在http中是通过Flusher实现，分批写入到http的writer，类似写入文件</p>

[Ref](https://stackoverflow.com/questions/26769626/send-a-chunked-http-response-from-a-go-server)

<br>
</br>

### 使用
1. 启动服务 
```shell
    go run simple.go
```

2. 用telnet请求(http会优化显示，用telnet直接打印完整的resp)
```
    telnet localhost 18080
```