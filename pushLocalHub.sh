# 本地打包并上传到内网hub，用于测试部署到容器的效果
docker kill fschedule
docker rm fschedule
docker rmi fschedule:latest

# 编译应用
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o fschedule-server -ldflags="-w -s" .
# 打包
docker build -t fschedule:latest --network=host .
# 发到内网
docker tag fschedule:latest hub.fsgit.cc/fschedule:dev
docker push hub.fsgit.cc/fschedule:dev
docker rmi hub.fsgit.cc/fschedule:dev
docker rmi fschedule:latest