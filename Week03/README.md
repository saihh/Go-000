学习笔记

启动两个goroutine分别启动服务和监听sigint信号，当有一个退出时两者都退出，并且等待http server优雅关闭后才终止程序。

启动程序后，ctr+c，输出如下

```
sigint notified
rcv stop and svr is about to shutdown
errgroup err sigint notified
svr truely shutdown
```