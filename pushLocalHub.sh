# 更新farseer-go框架
cd ../farseer-go && sh git-update.sh
# 更新FSchedule
cd ../FSchedule && git pull
# 生成路由
#fsctl -r
# 将忽略文件复制到上下文根目录中
#\cp .dockerignore ../
# 编译
docker build -t fschedule:latest --network=host -f ./Dockerfile ../
# 发到内网
docker tag fschedule:latest hub.fsgit.cc/fschedule:dev && docker push hub.fsgit.cc/fschedule:dev
docker rmi hub.fsgit.cc/fschedule:dev && docker rmi fschedule:latest

# docker
docker service rm fschedule
docker service create --name fschedule --replicas 1 -d --network=net \
--constraint node.role==worker \
--mount type=bind,src=/etc/localtime,dst=/etc/localtime \
-e "Database_default=\"DataType=mysql,PoolMaxSize=50,PoolMinSize=1,ConnectionString=root:123456@tcp(127.0.0.1:3306)/fschedule?charset=utf8&parseTime=True&loc=Local\"" \
-e "Redis_default=\"Server=127.0.0.1:6379,DB=14,Password=123456,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000\"" \
-l "traefik.http.routers.fschedule.rule=Host(\`fschedule.fsgit.cc\`)" \
-l "traefik.http.routers.fschedule.entrypoints=websecure" \
-l "traefik.http.routers.fschedule.tls=true" \
-l "traefik.http.services.fschedule.loadbalancer.server.port=8886" \
hub.fsgit.cc/fschedule:dev
