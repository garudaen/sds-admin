# SDS Admin - DNS 管理系统

一个现代化的 DNS 管理系统，采用 Golang (Gin) 后端和 Vue.js 3 前端构建。

[English](./README.md) | 简体中文

## 功能特性

- **域名管理**：添加、编辑、删除、启用/禁用域名
- **DNS 记录管理**：支持 A、AAAA、CNAME、MX、TXT 记录类型
- **智能 DNS 解析**：基于客户端 IP CIDR 的 A/AAAA/CNAME 记录路由
- **RESTful API**：完善的 Swagger API 文档支持
- **现代化 UI**：响应式 Web 界面，支持移动端

## 技术栈

### 后端
- **框架**：Gin (Golang)
- **ORM**：GORM
- **数据库**：MySQL
- **日志**：Logrus + Lumberjack 日志轮转
- **API 文档**：Swagger/OpenAPI

### 前端
- **框架**：Vue.js 3 (Composition API)
- **构建工具**：Vite
- **HTTP 客户端**：Axios
- **静态文件**：使用 Statik 内嵌

## 项目结构

```
sds-admin/
├── cmd/
│   └── sds-admin/          # 应用入口
│       └── main.go
├── configs/
│   ├── config-sample.yaml  # 配置示例
│   └── config.yaml         # 实际配置（已忽略）
├── internal/
│   ├── config/             # 配置加载
│   ├── database/           # 数据库连接和迁移
│   ├── dto/                # 数据传输对象
│   ├── handler/            # HTTP 处理器
│   ├── logger/             # 日志配置
│   ├── models/             # 数据库模型
│   ├── router/             # 路由定义
│   ├── service/            # 业务逻辑
│   └── static/             # 内嵌静态文件（生成）
├── pkg/                    # 公共库
├── fe/                     # 前端源码
│   ├── src/
│   │   ├── App.vue
│   │   └── components/
│   │       ├── AdminLayout.vue
│   │       ├── DomainManagement.vue
│   │       └── RecordManagement.vue
│   └── package.json
├── docs/                   # Swagger 文档（生成）
├── logs/                   # 日志文件（已忽略）
├── pub/                    # 构建的前端文件（已忽略）
├── bin/                    # 编译后的二进制文件（已忽略）
└── Makefile
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### 安装步骤

1. **克隆仓库**
   ```bash
   git clone <repository-url>
   cd sds-admin
   ```

2. **配置应用**
   ```bash
   cp configs/config-sample.yaml configs/config.yaml
   # 编辑 configs/config.yaml，填入数据库连接信息
   ```

3. **创建数据库**
   ```sql
   CREATE DATABASE sds_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

4. **构建并运行**
   ```bash
   make build-all
   ./bin/sds-admin --config configs/config.yaml
   ```

   或直接运行：
   ```bash
   make run
   ```

5. **访问应用**
   - Web 界面：http://localhost:8080
   - Swagger API 文档：http://localhost:8080/swagger/index.html

## API 端点

### 域名管理

| 方法 | 路径 | 描述 |
|--------|------|-------------|
| GET | `/api/v1/domains` | 获取域名列表 |
| GET | `/api/v1/domains/:id` | 根据 ID 获取域名 |
| POST | `/api/v1/domains` | 创建域名 |
| PUT | `/api/v1/domains/:id` | 更新域名 |
| DELETE | `/api/v1/domains/:id` | 删除域名 |
| POST | `/api/v1/domains/:id/disable` | 禁用域名 |
| POST | `/api/v1/domains/:id/enable` | 启用域名 |

### 记录管理

| 方法 | 路径 | 描述 |
|--------|------|-------------|
| GET | `/api/v1/domains/:id/records` | 获取域名的记录列表 |
| GET | `/api/v1/domains/:id/records/:recordId` | 根据 ID 获取记录 |
| POST | `/api/v1/domains/:id/records` | 创建记录 |
| PUT | `/api/v1/domains/:id/records/:recordId` | 更新记录 |
| DELETE | `/api/v1/domains/:id/records/:recordId` | 删除记录 |
| POST | `/api/v1/domains/:id/records/:recordId/disable` | 禁用记录 |
| POST | `/api/v1/domains/:id/records/:recordId/enable` | 启用记录 |

### 记录类型

| 方法 | 路径 | 描述 |
|--------|------|-------------|
| GET | `/api/v1/record-types` | 获取所有记录类型 |
| GET | `/api/v1/record-types/:id` | 根据 ID 获取记录类型 |

## DNS 记录类型说明

### A 记录
- 将域名映射到 IPv4 地址
- 支持通过 CIDR 实现多值智能解析
- 必须至少有一个默认值

### AAAA 记录
- 将域名映射到 IPv6 地址
- 支持通过 CIDR 实现多值智能解析
- 必须至少有一个默认值

### CNAME 记录
- 创建指向另一个域名的别名
- 同一主机记录不能与其他类型共存
- 自动补全末尾的点号

### MX 记录
- 邮件交换记录
- 每个记录值有独立的优先级（0-65535）
- 优先级数字越小，优先级越高

### TXT 记录
- 文本记录，用于多种用途
- 简单的值存储

## 开发指南

### Make 命令

```bash
make help          # 显示所有可用命令
make deps          # 安装 Go 依赖
make build         # 仅构建后端
make build-all     # 构建前端和后端
make run           # 运行应用
make clean         # 清理构建产物
make test          # 运行测试
make fmt           # 格式化代码
make lint          # 运行代码检查
make swagger       # 生成 Swagger 文档
make docker-build  # 构建 Docker 镜像
make docker-run    # 运行 Docker 容器
```

### 前端开发

```bash
cd fe
npm install
npm run dev     # 开发服务器（热重载）
npm run build   # 生产构建
```

### 数据库迁移

首次运行时自动创建数据表：
- `domains` - 域名表
- `record_types` - DNS 记录类型表（A、AAAA、CNAME、MX、TXT）
- `records` - DNS 记录表
- `record_values` - 记录值表（支持 CIDR）

## 配置说明

详见 `configs/config-sample.yaml`：

```yaml
server:
  host: 0.0.0.0
  port: 8080
  mode: release          # debug, release, test

database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: sds_admin

log:
  level: info            # debug, info, warn, error
  format: json           # json, text
  output: logs/sds-admin.log
  max_size: 100          # MB
  max_backups: 3
  max_age: 7             # days
  compress: true

swagger:
  enabled: true
```

## 开源协议

**WTFPL (Do What The Fuck You Want To Public License)**

本项目完全由 AI (Qoder) 生成。你可以用它做任何想做的事——无任何限制，无任何担保。

```
            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                    Version 2, December 2004

 Copyright (C) 2024 AI Generated

 Everyone is permitted to copy and distribute verbatim or modified
 copies of this license document, and changing it is allowed as long
 as the name is changed.

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

  0. You just DO WHAT THE FUCK YOU WANT TO.
```
