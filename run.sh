docker kill fschedule
docker rm fschedule
docker run --name fschedule -d --network=net \
-l "traefik.http.routers.vip8.rule=Host(\`fschedule.fsgit.cc\`)" \
-l "traefik.http.routers.vip8.entrypoints=web" \
-l "traefik.http.services.vip8.loadbalancer.server.port=8886" \
fschedule:latest