# download-geektime-video
此工具为下载极客时间已购课程方便离线观看
参考：https://github.com/domliang/geektime-dl
### 环境要求
```
（1）golang环境
（2）已安装 ffmpeg：brew install ffmpeg

```
### 使用

#### 1.进入项目目录，视频下载地址即 download-geektime-video.go 所在的目录

#### 2.编辑配置文件
```
vim config.json
```
配置文件示例
```
{
  "cid": "98",
  "_ga": "GA1.2.1666006123.1560340327",
  "_gid": "GA1.2.1666006123.1560340327",
  "GCID": "7d63098-852651f-85359ab-348d01c",
  "GCESS": "BAYEXcpaOwEEBoYRAAUEAAAAAAIEOI0HXQgBAwsCBAAMAQEKBAAAAAAEBAAvDQAHBOioBL0JAQEDBDiNB10-"
}
```
##### cid获取
浏览器打开并登录极客时间打开某一课程
```
https://time.geekbang.org/course/detail/168-68568
```
其中168是cid

##### cookie获取
拿到 .time.geekbang.org 下面的4个cookie：_ga，_gid，GCID，GCESS

#### 3.执行下载
```
go run download-geektime-video.go
```
