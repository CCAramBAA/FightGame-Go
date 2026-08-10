# 生成开发用 SSL 证书
# 方式一：使用 OpenSSL（推荐）
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout docker/nginx/ssl/server.key \
    -out docker/nginx/ssl/server.crt \
    -subj "/C=CN/ST=Beijing/L=Beijing/O=FightGame/CN=localhost"

# 方式二：使用 Docker 生成
docker run --rm -v $(pwd)/docker/nginx/ssl:/certs alpine/openssl \
    req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout /certs/server.key \
    -out /certs/server.crt \
    -subj "/CN=localhost"
