
```
docker run --name=systemstat -d -p 30000:30000 --net=mynetwork -e PREFIX=sysinfo -v /:/host:ro --restart=unless-stopped --privileged linuxshots/systemstat:1.0.0-arm64
```