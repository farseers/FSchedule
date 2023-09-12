docker kill fschedule
docker rm fschedule
docker run --name fschedule -d --restart=always \
-e Database_default="DataType=mysql,PoolMaxSize=50,PoolMinSize=1,ConnectionString=root:123456@tcp(127.0.0.1:3306)/fschedule?charset=utf8&parseTime=True&loc=Local" \
-e Redis_default="Server=127.0.0.1:6379,DB=14,Password=123456,ConnectTimeout=600000,SyncTimeout=10000,ResponseTimeout=10000" \
steden88/fschedule:latest
docker logs fschedule