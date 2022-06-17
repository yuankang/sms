### 需要解决的问题
1 播放发送数据，要修改fmt值和修改协议头
2 多个播放在一个协程里，如何避免互相阻塞
golang网络通信超时设置
https://www.cnblogs.com/lanyangsh/p/10852755.html

### sms  
Stream Media Server  

### 流媒体服务器
https://m7s.live/
https://j.m7s.live/player.html
https://www.freecodecamp.org/
https://chinese.freecodecamp.org/
https://www.topgoer.com/

### 编译命令  
go build -o sms main.go http.go rtmp.go  

### 支持的协议  
rtmp -> sms -> rtmp  

### rtmp参考资料
https://blog.csdn.net/lightfish_zhang/article/details/88681828
https://www.cnblogs.com/jimodetiantang/p/8974075.html

### 测试时使用到的命令
go build -o slivehlsupload_mac main.go http.go sqlite.go upload.go mqtt.go stream.go hls.go s3.go
ffmpeg -i GSPm8n46bmUqT-cf6dbn1ehZ_20220105103534_42840.ts -ss 0 -frames:v 1 -f image2 -y test.jpg
ffmpeg -re -stream_loop -1 -i fruit.mp4 -vcodec copy -acodec copy -f flv rtmp://127.0.0.1/live/cctv1
ffmpeg -re -stream_loop -1 -i fruit.mp4 -vcodec copy -acodec copy -f flv rtmp://172.20.25.20/SP63nBbfmlbW/GSP63nBbfmlbW-fnMebne7hV
ffplay rtmp://172.20.25.20/SP63nBbfmlbW/GSP63nBbfmlbW-fnMebne7hU
ffplay http://172.20.25.20:8082/SP63nBbfmlbW/GSP63nBbfmlbW-fnMebne7hU.flv?mediaServerIp=172.20.25.20&codeType=H264
ffplay http://172.20.25.20:8082/SP63nBbfmlbW/GSP63nBbfmlbW-fnMebne7hV.m3u8

### sms待解决问题
1 rtmp收流，rtmp/flv/hls播放
2 支持gateway推流，检查ts生成是否正常
3 支持obs推流，检查ts生成是否正常
4 支持h265
5 支持http截图请求
6 支持flv播放加密
7 停止推流程序崩溃
8 内存使用优化
9 流量统计和上报
10 推流鉴权
11 流状态变更上报
12 支持http获取推流的流id
13 支持http获取每路流的播放个数
14 rtmp推流
15 map要使用sync.Map
	https://pkg.go.dev/sync@go1.17.8
16 使用mqtt发布tsinfo
17 rtp(udp)收流
18 rtp(tcp)收流
19 rtp(udp)推流
20 rtp(tcp)推流
21 rtcp客户端实现
22 rtcp服务端实现
23 音频重采样，出固定采样率
24 rtsp支持
25 webrtc支持
26 websocket支持
