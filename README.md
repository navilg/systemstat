
```
docker run --name=systemstat -d -p 30000:30000 --net=mynetwork -e PREFIX=sysinfo -v /:/host:ro --restart=unless-stopped --privileged linuxshots/systemstat:1.0.0-arm64
```

# Put behind Nginx

```
location /sysinfo {
    proxy_set_header Host $host;                                         
    proxy_set_header X-Forwarded-Scheme $scheme;                         
    proxy_set_header X-Forwarded-Proto $scheme;                          
    proxy_set_header X-Forwarded-For $remote_addr;     
    proxy_set_header X-Real-IP $remote_addr;                             
    proxy_pass http://systemstat:30000;
}
```

Add Basic Auth password in Nginx

- If your system is debian-based install apache2-utils

```
sudo apt install apache2-utils
```

- Generate a .htpasswd file

```
sudo htpasswd -c .htpasswd admin
```

- Move .hdpasswd file to Nginx container

```
docker cp .htpasswd nginx:/etc/nginx/.htpasswd
```

- Add basic auth in proxy

e.g.

```
location /sysinfo {
    proxy_set_header Host $host;                                         
    proxy_set_header X-Forwarded-Scheme $scheme;                         
    proxy_set_header X-Forwarded-Proto $scheme;                          
    proxy_set_header X-Forwarded-For $remote_addr;     
    proxy_set_header X-Real-IP $remote_addr;                             
    proxy_pass http://systemstat:30000;
    auth_basic           "Administrator’s Area";
    auth_basic_user_file /etc/nginx/.htpasswd; 
}
```

- Reload Nginx

```
nginx -s reload
```