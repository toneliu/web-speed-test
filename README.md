# 网络带宽测速系统

一个基于 Go 语言的网络带宽测速系统，支持多单位管理、测速记录、拓扑可视化等功能。

## 功能特性

- **客户端测速**：用户选择所属单位后点击测速，自动测量下行带宽并上报结果，支持切换单位（Mbps/MB/s/Kbps）
- **自动上报与历史记录**：每次测速自动保存。管理员可查看所有单位的测速历史（表格 + 趋势折线图）；各单位有自己的账号，登录后只能查看自己单位的测速历史记录
- **多单位与单位账号管理**：管理员可增删改查单位，并可为每个单位创建/重置登录账号（用户名 + 密码）
- **拓扑可视化**：星型拓扑图展示服务器到各单位的链路带宽，支持管理员手动添加单位间自定义链路
- **单文件部署**：生成一个可执行文件，运行即用（前端静态资源用 go:embed 嵌入，数据库用 SQLite）
- **GitHub Actions 自动构建**：自动编译 Linux 和 Windows 版本并发布到 Release

## 默认账号

| 角色     | 用户名   | 密码       | 说明         |
|----------|----------|------------|--------------|
| 管理员   | admin    | admin123   | 系统管理员   |
| 单位A用户| unita    | unita123   | 单位A账号    |
| 单位B用户| unitb    | unitb123   | 单位B账号    |
| 单位C用户| unitc    | unitc123   | 单位C账号    |

## 快速开始

### 从 Release 下载

1. 访问 [Releases 页面](https://github.com/your-username/speedtest/releases)
2. 下载对应系统的可执行文件：
   - Linux (386): `speedtest-linux-386`
   - Linux (AMD64): `speedtest-linux-amd64`
   - Windows (386): `speedtest-windows-386.exe`
   - Windows (AMD64): `speedtest-windows-amd64.exe`
3. 运行程序：
   ```bash
   # Linux
   chmod +x speedtest-linux-amd64
   ./speedtest-linux-amd64

   # Windows (PowerShell)
   .\speedtest-windows-amd64.exe
   ```

### 本地构建

```bash
# 克隆仓库
git clone https://github.com/your-username/speedtest.git
cd speedtest

# 安装依赖
go mod tidy

# 运行开发服务器
go run main.go

# 构建可执行文件
go build -o speedtest .

# 使用脚本构建所有平台版本
./build.sh
```

### 访问应用

- 打开浏览器访问：http://localhost:8080
- 使用默认账号登录

### 开始测速

- 选择单位
- 点击"开始测速"
- 查看结果和历史记录

## 项目结构

```
speedtest/
├── main.go              # 主程序入口
├── go.mod               # Go 依赖管理
├── go.sum               # Go 依赖锁定
├── build.sh             # 本地构建脚本
├── .gitignore           # Git 忽略文件
├── README.md            # 项目文档
├── .github/
│   └── workflows/
│       └── release.yml  # GitHub Actions 配置
├── frontend/
│   └── index.html       # 单页面应用
└── pkg/
    ├── models/          # 数据模型
    ├── database/        # 数据库操作
    ├── middleware/      # 认证中间件
    └── handlers/        # API 处理函数
```

## 开发

```bash
# 安装依赖
go mod tidy

# 运行开发服务器
go run main.go

# 构建可执行文件
go build -o speedtest .

# 构建所有平台版本
./build.sh
```

## GitHub Actions 自动构建和发布

项目配置了 GitHub Actions 自动构建和发布流程：

### 触发方式

1. **推送 Tag**：当推送格式为 `v*` 的 tag 时（如 `v1.0.0`），自动触发构建和发布
2. **手动触发**：在 GitHub 仓库的 Actions 页面手动触发

### 发布流程

- 自动构建 Linux (386/AMD64) 和 Windows (386/AMD64) 版本
- 将构建产物上传到 Artifacts
- 自动创建 GitHub Release 并附加所有编译好的文件

### 创建发布

```bash
# 创建 tag
git tag -a v1.0.0 -m "Release version 1.0.0"

# 推送 tag
git push origin v1.0.0
```

## 技术栈

- **后端**：Go + Gin + GORM + SQLite + JWT
- **前端**：原生 HTML/JS + Chart.js
- **打包**：go:embed 实现单文件部署
- **CI/CD**：GitHub Actions
