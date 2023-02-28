ver='2.0.0'
# 编译应用
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o fschedule-server -ldflags="-w -s" .
# 打包
docker build -t steden88/fschedule:${ver} --network=host .
docker tag steden88/fschedule:${ver} steden88/fschedule:latest

docker push steden88/fschedule:${ver}
docker push steden88/fschedule:latest