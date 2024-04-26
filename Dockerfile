FROM alpine:latest
WORKDIR /app
# 复制配置（没有配置需要注释掉）
COPY /FSchedule/farseer.yaml .
COPY /FSchedule/FSchedule .

#设置时区
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai    /etc/localtime

ENTRYPOINT ["./FSchedule"]