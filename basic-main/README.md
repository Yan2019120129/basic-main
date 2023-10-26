# basic
zfeng



#### Quasar
    1.创建方式 - 使用 vue3 + TypeScript + Vite
        yarn create quasar
        quasar dev -m ssr                       本地测试启动 ssr模式
        quasar dev                              本地测试启动 spa模式



#### 反响代理
server {
listen 80;
server_name api.beego.com;
    location / {
        proxy_pass http://127.0.0.1:9090;
        proxy_set_header Upgrade websocket;
        proxy_set_header Connection Upgrade;
    }
}

vue 伪静态
location / {
    try_files $uri $uri/ /index.html;
}

proxy_pass http://127.0.0.1:3001;
proxy_set_header Host $host:$server_port;
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header REMOTE-HOST $remote_addr;
add_header X-Cache $upstream_cache_status;
proxy_set_header X-Host $host:$server_port;
proxy_set_header X-Scheme $scheme;
proxy_connect_timeout 30s;
proxy_read_timeout 86400s;
proxy_send_timeout 30s;
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o src/adm.service admin/run.go 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o src/web.service home/run.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o src/web.service home/run.go


安装node
curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
apt-get install -y nodejs
sudo apt-get install yarn


设置世界标准时间
timedatectl set-timezone UTC