# 杀掉 air 进程，防止多次启动，杀掉之前运行的 air 进程
ps -ef | grep -w air | grep -v grep | sort -k 2rn | awk '{if (NR>1){print $2}}' | xargs kill -9
# 杀掉 serve 进程，这里的端口是 gin 运行的端口
lsof -i:8001 | grep serve | awk '{print $2}' | xargs kill -9
# 杀掉 dlv 进程，这里的端口是 dlv 运行的端口
lsof -i:2345 | grep dlv | awk '{print $2}' | xargs kill -9
# debug gin 项目
dlv debug --listen=:2345 --headless=true --api-version=2 --continue --accept-multiclient --output=./tmp/serve ./admin/run.go
