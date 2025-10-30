# GitHub Pages 设置指南

## 📋 前提条件

- ✅ 已创建 `.github/workflows/hugo.yml` 文件
- ✅ Hugo 配置文件 `hugo.toml` 中的 `baseURL` 已正确设置

## 🔧 GitHub 仓库设置

### 1. 启用 GitHub Pages

1. 访问你的 GitHub 仓库：https://github.com/heyuuuu77/blog
2. 点击 **Settings** (设置)
3. 在左侧菜单找到 **Pages**
4. 在 **Build and deployment** 部分：
   - **Source**: 选择 **GitHub Actions**
   - ⚠️ 不要选择 "Deploy from a branch"

### 2. 配置权限（如果需要）

如果部署失败，可能需要配置权限：

1. 在仓库的 **Settings** → **Actions** → **General**
2. 滚动到 **Workflow permissions**
3. 选择 **Read and write permissions**
4. 勾选 **Allow GitHub Actions to create and approve pull requests**
5. 点击 **Save**

## 🚀 部署流程

### 自动部署

每次推送到 `main` 分支时，GitHub Actions 会自动：

1. ✅ 检出代码
2. ✅ 安装 Hugo
3. ✅ 构建静态网站
4. ✅ 部署到 GitHub Pages

### 手动触发

1. 访问 https://github.com/heyuuuu77/blog/actions
2. 选择 **Deploy Hugo site to Pages** workflow
3. 点击 **Run workflow**
4. 选择分支 `main`
5. 点击 **Run workflow**

## 📝 日常使用

### 发布新文章

```bash
# 1. 创建新文章
hugo new posts/my-new-post.md

# 2. 编写内容
vim content/posts/my-new-post.md

# 3. 本地预览
hugo serve

# 4. 提交并推送（自动触发部署）
git add .
git commit -m "Add: my new post"
git push origin main
```

### 查看部署状态

1. 推送后访问：https://github.com/heyuuuu77/blog/actions
2. 查看最新的 workflow 运行状态
3. 点击查看详细日志

部署通常需要 1-2 分钟完成。

## 🌐 访问你的博客

部署成功后，访问：
- **标准地址**: https://heyuuuu77.github.io/blog/
- 或者 **自定义域名**（如果已配置）

⚠️ **注意**: 如果你的仓库名是 `heyuuuu77.github.io`，则直接访问 https://heyuuuu77.github.io/

## 🐛 常见问题

### 问题 1: Workflow 运行失败

**解决方案**：
1. 检查 Actions 权限设置（见上文）
2. 确保 `hugo.toml` 中的配置正确
3. 查看 workflow 日志获取详细错误信息

### 问题 2: 页面 404

**解决方案**：
1. 确认 GitHub Pages Source 设置为 **GitHub Actions**
2. 检查 `hugo.toml` 中的 `baseURL` 是否正确：
   ```toml
   baseURL = "https://heyuuuu77.github.io/blog/"
   ```
3. 如果仓库名是 `heyuuuu77.github.io`，baseURL 应该是：
   ```toml
   baseURL = "https://heyuuuu77.github.io/"
   ```

### 问题 3: 样式加载失败

**解决方案**：
- 检查 `baseURL` 配置
- 清除浏览器缓存
- 等待几分钟让 CDN 刷新

### 问题 4: Workflow 权限错误

错误信息类似：
```
Error: Resource not accessible by integration
```

**解决方案**：
1. 进入 Settings → Actions → General
2. 设置 Workflow permissions 为 "Read and write permissions"

## 📊 Workflow 说明

### 触发条件

- ✅ 推送到 `main` 分支
- ✅ 手动触发（workflow_dispatch）

### 构建步骤

1. **安装 Hugo CLI** (extended 版本)
2. **安装 Dart Sass** (处理 SCSS)
3. **检出代码** (包括子模块)
4. **配置 Pages** (设置 baseURL)
5. **构建网站** (`hugo --gc --minify`)
6. **上传产物** (public 目录)
7. **部署** (deploy-pages action)

### Hugo 版本

当前配置使用 Hugo 0.128.0。如需更改版本，修改 `hugo.yml` 中的：

```yaml
env:
  HUGO_VERSION: 0.128.0  # 修改为所需版本
```

## 🎯 下一步

1. ✅ 推送 workflow 文件到 GitHub
2. ✅ 在 GitHub 仓库设置中启用 Pages (Source: GitHub Actions)
3. ✅ 推送一次代码触发首次部署
4. ✅ 等待部署完成（1-2 分钟）
5. ✅ 访问你的博客网址

## 📚 参考文档

- [Hugo 官方文档](https://gohugo.io/documentation/)
- [GitHub Pages 文档](https://docs.github.com/en/pages)
- [GitHub Actions 文档](https://docs.github.com/en/actions)

