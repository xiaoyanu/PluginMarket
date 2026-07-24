# PluginMarket

一个基于 **Nuxt4 + Go** 构建的插件分享与管理平台。

PluginMarket 面向插件作者和使用者，提供插件发布、审核、搜索、下载、评论、回复、收藏以及站内通知等功能，同时提供管理员后台用于管理插件、用户、评论和系统通知。

## 功能特性

### 用户侧

- 用户注册、登录、验证码和密码找回
- 用户资料、头像、邮箱和称号管理
- 插件列表、热门插件和最新插件浏览
- 插件搜索及按标签、框架、作者筛选
- 插件详情、下载地址和提取码展示
- 插件收藏
- 评论和多级回复
- 评论分页：每页 20 条顶级评论，子评论随所属评论展示
- 评论通知深链接：从邮件或站内通知打开后，自动定位并高亮目标评论
- 子评论自动展开和精确定位
- 站内通知未读数、未读优先排序和一键已读

### 管理侧

- 插件审核、发布、编辑和删除
- 用户列表、筛选、分页和管理员操作
- 标签、框架和称号管理
- 网站基础设置和邮件模板设置
- 站内通知创建、编辑、筛选和删除
- 支持系统通知、自定义通知、新增评论、评论回复、审核通过和审核拒绝等通知类型
- 通知列表使用统一的类型标签和状态展示

## 技术栈

### 前端

- [Vue 3](https://vuejs.org/)
- [Nuxt 4](https://nuxt.com/)
- [TypeScript](https://www.typescriptlang.org/)
- [Element Plus](https://element-plus.org/)
- [Tailwind CSS](https://tailwindcss.com/)
- [SCSS](https://sass-lang.com/)
- [Pinia](https://pinia.vuejs.org/)
- [md-editor-v3](https://github.com/imzbf/md-editor-v3)
- [Phosphor Icons](https://phosphoricons.com/)

### 后端

- [Go](https://go.dev/) 1.25+
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [MySQL](https://www.mysql.com/) 5.6+
- JWT 身份认证
- YAML 配置
- SMTP 邮件发送

## 项目结构

```text
PluginMarket/
├── server/                 # Go 后端服务
│   ├── config.yaml         # 后端运行配置示例
│   ├── controller/         # HTTP 控制器和业务流程
│   ├── database/           # 数据库连接与迁移
│   ├── middleware/         # 认证、管理员和跨域中间件
│   ├── model/              # GORM 数据模型
│   ├── repository/         # 数据访问层
│   ├── router/             # API 路由
│   ├── uploads/            # 上传文件目录
│   └── main.go             # 服务入口
├── web/                    # Nuxt 前端应用
│   ├── app/components/     # Vue 组件
│   ├── app/composables/    # API 请求和组合式逻辑
│   ├── app/pages/          # 页面路由
│   ├── app/utils/          # 前端工具函数
│   ├── nuxt.config.ts      # Nuxt 配置
│   └── package.json        # 前端依赖和脚本
├── LICENSE
└── README.md
```

## 环境要求

- Go 1.25 或更高版本
- Node.js 20+，建议使用当前 LTS 版本
- pnpm 9+
- MySQL 5.6 或更高版本
- Git

## 快速开始

### 1. 获取项目

```bash
git clone https://github.com/xiaoyanu/PluginMarket.git
cd PluginMarket
```

### 2. 准备数据库

创建数据库并导入初始结构：

```sql
CREATE DATABASE tpm CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

```bash
mysql -h127.0.0.1 -P3306 -uroot -p tpm < server/db/pm.sql
```

根据本地环境修改 `server/config.yaml` 中的数据库连接配置。不要将生产密码、SMTP 密码或 JWT 密钥提交到 Git。

### 3. 启动后端

```bash
cd server
go mod download
go run .
```

默认监听：

```text
http://127.0.0.1:8855
```

### 4. 安装并启动前端

新开一个终端：

```bash
cd web
pnpm install
pnpm dev --host 0.0.0.0 --port 3000
```

默认访问：

```text
http://127.0.0.1:3000
```

前端 API 地址和资源地址由 `web/nuxt.config.ts` 配置。部署到其他环境时，请根据实际网络拓扑修改，不要把本机地址作为生产配置提交。

### 5. 检查服务

```bash
curl -I http://127.0.0.1:3000/
curl http://127.0.0.1:8855/api/setting/public
ss -ltnp | grep -E ':(3000|8855|3306)\\b'
```

预期结果：

- 前端 `3000` 端口返回 HTTP 200
- 后端公开设置接口返回 `code: 200`
- MySQL 监听 `3306`

## 配置说明

后端主要配置位于 `server/config.yaml`：

| 配置项 | 说明                |
| --- |-------------------|
| `server.port` | 后端监听端口，默认 `8855`  |
| `server.webUrl` | 邮件链接使用的站点公开地址     |
| `database.*` | MySQL 连接信息        |
| `jwt.secret` | JWT 签名密钥，生产环境必须替换 |
| `jwt.expire` | JWT 有效期，单位为小时     |
| `uploads.path` | 上传文件根目录           |
| `uploads.maxSize` | 上传文件大小限制          |

生产环境建议：

- 使用环境隔离的数据库账号，不使用默认 `root` 账号
- 使用随机且长度足够的 JWT 密钥
- 使用 HTTPS 的公开站点地址
- 将上传目录放在持久化存储中
- 不提交包含真实凭据的配置文件

## 开发命令

### 后端

```bash
cd server
go test ./...
go fmt ./...
```

### 前端

```bash
cd web
pnpm dev
pnpm build
pnpm generate
pnpm preview
```

前端轻量测试文件使用 Node.js 原生测试运行器，例如：

```bash
node --test app/utils/*.test.ts
```

## 数据和上传文件

- 服务启动时会执行必要的迁移和初始化逻辑。
- `server/uploads/` 用于保存上传的图片、头像、标题和框架资源。
- 测试上传数据属于项目运行所需数据，清理前请确认不会影响本地验证。
- 生产部署时应使用独立的数据目录和备份策略。

## 贡献指南

1. 从最新的 `master` 分支创建功能分支。
2. 保持修改范围聚焦，避免混入无关格式化或重构。
3. 为新的业务行为补充回归测试。
4. 提交前运行后端测试和前端构建。
5. 使用清晰的 Conventional Commits 风格提交信息，例如：

```text
feat: add plugin comment pagination
fix: prioritize unread notifications
docs: improve project README
```

6. 提交 Pull Request 时说明：
   - 修改内容
   - 影响范围
   - 测试结果
   - 是否需要数据库迁移或配置变更

## 许可证

本项目使用 [MIT License](LICENSE)。
