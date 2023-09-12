docker service rm fschedule
docker service create --name fschedule --replicas 3 -d --network=net \
--constraint node.role==worker \
-l "traefik.http.routers.fschedule.rule=Host(\`fschedule.fsgit.cc\`)" \
-l "traefik.http.routers.fschedule.entrypoints=websecure" \
-l "traefik.http.routers.fschedule.tls=false" \
-l "traefik.http.services.fschedule.loadbalancer.server.port=8886" \
steden88/fschedule:latest