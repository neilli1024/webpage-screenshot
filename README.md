# 服务器渲染网页生成图片服务


作者：jessezhang007007 <jessezhang007007@gmail.com>

github地址：https://github.com/jessezhang007007/webpage-screenshot

Docker Hub 地址：https://hub.docker.com/r/jessezhang007007/webpage-screenshot/


## 一、说明：

服务器渲染网页生成图片服务，使用了chromedp运行chrome渲染

## 二、安装

### 方式一：使用Docker运行

```
docker run -d -p 80:8082 --name webpage-screenshot --restart=always jessezhang007007/webpage-screenshot
```


### 方式二：本地golang运行

```bash
go run app.go
```


## 三、访问服务
### 1. 传入url方式访问网页，使用GET方式访问：

```
http://localhost:8082/webpage/image?url=https%3A%2F%2Fwww.baidu.com
```

### 2. 直接传入html内容渲染，使用POST方式访问：

```
curl -X POST \
  http://localhost:8082/webpage/image \
  -o webpage-image.png \
  -d html=%3Chtml%3E%3Cbody%3E%E8%BF%99%E6%98%AFhtml%E5%86%85%E5%AE%B9%3C%2Fbody%3E%3C%2Fhtml%3E
```
