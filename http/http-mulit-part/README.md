### 基于http的multipart/form-data上传大文件
<br>
</br>
主要以下Ref的demo:
1. server
  解析mulit-form格式数据的simple server
2.client
  2.1 以mulit-form上传文件（实际server收到的还是stream-oct格式，如果想改可自定义multipart的CreateFormFile）
  2.2 基于io.Pipe()[ref](https://juejin.cn/post/7032933119992791053)实现大文件上传
    实际上传发现不仅可以减少app内存和带宽，还能减少传输时间
<br>
</br>



[Ref](https://tonybai.com/2021/01/16/upload-and-download-file-using-multipart-form-over-http/)

