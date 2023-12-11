# 更新farseer-go框架
cd ../farseer-go && sh git-update.sh
# 更新FSchedule
cd ../FSchedule && git pull
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
-e "Database_default=DataType=mysql,PoolMaxSize=25,PoolMinSize=1,ConnectionString=root:steden@123@tcp(192.168.1.8:3306)/fschedule?charset=utf8&parseTime=True&loc=Local" \
-e "Redis_default=Server=192.168.1.8:6379,DB=14,Password=steden@123,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000" \
-l "traefik.http.routers.fschedule.rule=Host(\`fschedule.fsgit.cc\`)" \
-l "traefik.http.routers.fschedule.entrypoints=websecure" \
-l "traefik.http.routers.fschedule.tls=true" \
-l "traefik.http.services.fschedule.loadbalancer.server.port=8886" \
hub.fsgit.cc/fschedule:dev