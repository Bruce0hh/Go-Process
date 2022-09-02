# go rpc

> 该项目所有内容均摘自https://geektutu.com/

## 处理超时

### 客户端处理超时

1. 与服务端建立连接导致的超时
2. 发送请求到服务端，写报文导致的超时
3. 等待服务端处理，等待处理导致的超时
4. 从服务端接收响应，读报文导致的超时

### 服务端处理超时

1. 读取客户端的请求报文，读报文导致的超时
2. 发送相应报文，写报文导致的超时
3. 调用映射服务的方法，处理报文导致的超时