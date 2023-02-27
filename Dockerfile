FROM alpine:latest AS base
WORKDIR /app
COPY ./fschedule-server .
# 复制配置
COPY ./farseer.yaml .
# 复制视图
#COPY ./views ./views
# 复制静态资源
#COPY ./wwwroot ./wwwroot

#设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai    /etc/localtime

ENTRYPOINT ["./fschedule-server"]