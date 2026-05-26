# 🚀 GitHub 仓库创建指南

由于 API token 权限限制，需要手动在 GitHub 网站上创建仓库。请按照以下步骤操作：

## 步骤 1: 创建 GitHub 仓库

1. **打开 GitHub 网站**
   访问 https://github.com 并登录您的账户

2. **创建新仓库**
   - 点击右上角的 **+** 按钮
   - 选择 **"New repository"**

3. **填写仓库信息**
   - **Repository name**: `web-speed-test`
   - **Description**: `Web Speed Test - 网络测速工具`
   - **选择 Public** (公共仓库)
   - **不要勾选** "Add a README file"
   - **不要勾选** "Add .gitignore"
   - **不要勾选** "Choose a license"

4. **点击 "Create repository"**

## 步骤 2: 连接本地仓库

创建完成后，GitHub 会显示新仓库页面。您会看到一个 HTTPS URL，类似于：
```
https://github.com/toneliu/web-speed-test.git
```

现在在终端中运行以下命令（假设您还在 web-speed-test 目录）：

```bash
# 添加远端仓库（将 URL 替换为您实际的仓库 URL）
git remote add origin https://github.com/toneliu/web-speed-test.git

# 重命名分支为 main
git branch -M main

# 推送代码
git push -u origin main
```

## 或者使用便捷脚本

我们为您准备了一个自动设置脚本：

```bash
cd /workspace/web-speed-test
bash setup.sh
```

脚本会引导您完成所有步骤。

## 步骤 3: 验证

推送成功后，您应该会看到：
```
To https://github.com/toneliu/web-speed-test.git
 * [new branch]      main -> main
Branch 'main' set up to track 'origin/main'.
```

现在您可以：
- 🌐 访问您的 GitHub 仓库: `https://github.com/toneliu/web-speed-test`
- 📦 查看代码
- 🎨 启用 GitHub Pages 进行在线预览

## 启用 GitHub Pages（可选）

1. 进入仓库设置 (Settings)
2. 滚动到 "GitHub Pages" 部分
3. Source 选择 `main` 分支和 `/(root)` 根目录
4. 点击 Save
5. 您的测速工具将可在 `https://toneliu.github.io/web-speed-test` 访问

## 项目文件说明

- `README.md` - 项目说明文档
- `index.html` - 主页面（包含测速功能）
- `setup.sh` - 自动设置脚本

## 本地测试

在推送之前，您可以在本地测试：

```bash
# 在浏览器中直接打开
open index.html  # macOS
xdg-open index.html  # Linux
start index.html  # Windows

# 或者使用本地服务器
python -m http.server 8000
# 然后访问 http://localhost:8000
```

## 🎉 恭喜！

您的 Web Speed Test 项目已准备就绪！
