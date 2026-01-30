# Homebox 部署指南

## 项目概述

**Homebox** 是一个协作标签和空间管理应用

**技术栈:**
- Go 1.25 + Gin 框架
- MySQL 8.0
- Redis 7.0
- JWT 认证

---

## 一、前置准备

### 1.1 环境要求

| 组件 | 要求 |
|------|------|
| CPU | 2核+ |
| 内存 | 4GB+ |
| 存储 | 20GB SSD |
| 系统 | Linux/macOS/Windows |

### 1.2 软件依赖

```bash
# Docker
docker --version    # 20.10+
docker-compose --version    # 2.0+
```

---

## 二、快速部署 (Docker Compose)

### 2.1 创建 .env 文件

```bash
cp .env.example .env
nano .env
```

```bash
# MySQL 容器配置
MYSQL_ROOT_PASSWORD=your-root-password
MYSQL_DATABASE=homebox
MYSQL_USER=homebox
MYSQL_PASSWORD=your-db-password

# Redis 容器配置
REDIS_PASSWORD=your-redis-password

# 应用配置
JWT_SECRET=your-jwt-secret-at-least-32-chars
SERVER_MODE=release
```

### 2.2 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: homebox_mysql
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: ${MYSQL_DATABASE}
      MYSQL_USER: ${MYSQL_USER}
      MYSQL_PASSWORD: ${MYSQL_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - homebox_network
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: homebox_redis
    command: redis-server --requirepass ${REDIS_PASSWORD}
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - homebox_network
    restart: unless-stopped

  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: homebox_app
    ports:
      - "8080:8080"
    environment:
      # 数据库连接（固定值，连接 docker-compose 网络内的服务）
      DATABASE_HOST: mysql
      DATABASE_PORT: 3306
      DATABASE_DBNAME: ${MYSQL_DATABASE}
      DATABASE_USER: ${MYSQL_USER}
      DATABASE_PASSWORD: ${MYSQL_PASSWORD}
      # Redis 连接
      REDIS_HOST: redis
      REDIS_PORT: 6379
      REDIS_PASSWORD: ${REDIS_PASSWORD}
      REDIS_DB: 0
      # 应用配置
      JWT_SECRET: ${JWT_SECRET}
      JWT_EXPIRE_TIME: 7200
      JWT_REFRESH_EXPIRE_TIME: 604800
      SERVER_PORT: "8080"
      SERVER_MODE: ${SERVER_MODE}
    depends_on:
      - mysql
      - redis
    networks:
      - homebox_network
    restart: unless-stopped

volumes:
  mysql_data:
  redis_data:

networks:
  homebox_network:
    driver: bridge
```

### 2.3 启动服务

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f app
```

### 2.4 验证部署

```bash
# 检查服务状态
docker-compose ps

# 健康检查
curl http://localhost:8080/test

# API 测试
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"test","email":"test@example.com","password":"password123"}'
```

---

## 三、传统部署

### 3.1 编译项目

```bash
# 克隆项目
git clone https://your-repo/homebox.git
cd homebox

# 编译
go build -o homebox .
```

### 3.2 配置文件

编辑 `config/config.yml`:

```yaml
database:
  host: localhost
  port: 3306
  dbname: homebox
  user: homebox
  password: your-password

redis:
  host: localhost
  port: 6379
  password: your-redis-password
  db: 0

jwt:
  secret: your-jwt-secret
  expire_time: 7200
  refresh_expire_time: 604800

server:
  port: 8080
  mode: release
```

#### 环境变量覆盖

应用支持通过环境变量覆盖配置文件，环境变量命名规则：**配置键 `.` 替换为 `_`，全大写**。

| 配置键 | 环境变量 |
|--------|----------|
| database.host | `DATABASE_HOST` |
| database.port | `DATABASE_PORT` |
| database.user | `DATABASE_USER` |
| database.password | `DATABASE_PASSWORD` |
| database.dbname | `DATABASE_DBNAME` |
| redis.host | `REDIS_HOST` |
| redis.port | `REDIS_PORT` |
| redis.password | `REDIS_PASSWORD` |
| redis.db | `REDIS_DB` |
| jwt.secret | `JWT_SECRET` |
| jwt.expire_time | `JWT_EXPIRE_TIME` |
| jwt.refresh_expire_time | `JWT_REFRESH_EXPIRE_TIME` |
| server.port | `SERVER_PORT` |
| server.mode | `SERVER_MODE` |

### 3.3 启动服务

```bash
# 确保 MySQL 和 Redis 已启动

# 运行应用
./homebox
```

### 3.4 使用 Nginx 反向代理 (可选)

```nginx
upstream homebox {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://homebox;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

---

## 四、常用命令

```bash
# Docker Compose
docker-compose up -d          # 启动服务
docker-compose down           # 停止服务
docker-compose logs -f app    # 查看日志
docker-compose ps             # 服务状态
docker-compose restart app    # 重启应用

# 传统部署
./homebox                     # 启动服务
nohup ./homebox &             # 后台运行
```

---

## 五、端口说明

| 服务 | 端口 |
|------|------|
| 应用 | 8080 |
| MySQL | 3306 |
| Redis | 6379 |
