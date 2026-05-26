#!/bin/bash

echo "🌐 Web Speed Test - GitHub 仓库设置脚本"
echo "========================================"
echo ""

# 检查是否在正确的目录
if [ ! -d ".git" ]; then
    echo "❌ 错误: 请在 web-speed-test 目录中运行此脚本"
    echo "   运行: cd web-speed-test && bash setup.sh"
    exit 1
fi

# 显示说明
echo "📋 步骤说明:"
echo "1. 在浏览器中打开 GitHub 并登录"
echo "2. 点击右上角的 '+' 按钮 -> 'New repository'"
echo "3. Repository name 填写: web-speed-test"
echo "4. 选择 Public (公共仓库)"
echo "5. 不要勾选任何初始化选项"
echo "6. 点击 'Create repository'"
echo ""
echo "7. 创建完成后，复制仓库的 HTTPS URL"
echo "   例如: https://github.com/toneliu/web-speed-test.git"
echo ""
read -p "按 Enter 键继续..."

# 获取仓库 URL
echo ""
read -p "请粘贴您的 GitHub 仓库 HTTPS URL: " repo_url

if [ -z "$repo_url" ]; then
    echo "❌ 错误: URL 不能为空"
    exit 1
fi

# 添加远端仓库
echo ""
echo "🔗 添加远端仓库..."
git remote add origin "$repo_url"

# 推送代码
echo ""
echo "📤 推送代码到 GitHub..."
git branch -M main
git push -u origin main

echo ""
echo "✅ 完成！"
echo "您的项目已成功推送到: $repo_url"
echo ""
echo "🎉 恭喜！您可以访问您的 GitHub 仓库了！"
