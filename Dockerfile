# 注意，这里的构建上下文，是在git源代码的根目录
FROM golang:1.20.11-alpine AS build
# 设置github代理
ENV GOPROXY https://goproxy.cn,direct
# 进入到项目目录中
WORKDIR /src/FSchedule
# 复制go.mod文件
COPY ./FSchedule/go.mod .
# 下载依赖（支持docker缓存）
RUN go mod download
# 将源代码复制到此
COPY ./FSchedule .
# 删除go.work文件
#RUN rm -rf go.work
# 更新go.sum
RUN go mod tidy
# farseer项目
WORKDIR /src/farseer-go
COPY ./farseer-go .
# 进入到项目目录中
WORKDIR /src/FSchedule
# 编译
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/fschedule-server -ldflags="-w -s" .

FROM alpine:latest AS base
WORKDIR /app
COPY --from=build /app .
# 复制配置（没有配置需要注释掉）
COPY --from=build /src/FSchedule/farseer.yaml .
# 复制视图（没有视图需要注释掉）
#COPY --from=build /src/views ./views
# 复制静态资源（没有静态资源需要注释掉）
#COPY --from=build /src/wwwroot ./wwwroot

#设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai    /etc/localtime

ENTRYPOINT ["./fschedule-server"]

