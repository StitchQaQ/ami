# 文章按更新日期排序说明

## ✅ 已完成的配置

在 `hugo.toml` 中添加了以下配置：

```toml
# Front Matter 日期配置（优先使用更新时间排序）
[frontmatter]
date = ["lastmod", "date", "publishDate"]
lastmod = [":git", "lastmod", "date", "publishDate"]
publishDate = ["publishDate", "date"]
```

### 配置说明

这个配置让 Hugo 按以下优先级读取日期：

1. **`date` 字段**：优先读取 `lastmod`，如果没有则读取 `date`，最后才是 `publishDate`
2. **`lastmod` 字段**：
   - `:git` - 优先使用 Git 最后提交时间
   - `lastmod` - 如果有手动设置的 `lastmod` 使用它
   - `date` - 回退到创建日期
   - `publishDate` - 最后才用发布日期

## 📝 如何使用

### 方法 1：手动设置更新时间（推荐）

在文章的 Front Matter 中添加 `lastmod` 字段：

```markdown
---
title: "文章标题"
date: 2025-01-15T10:00:00+08:00
lastmod: 2025-02-20T15:30:00+08:00  # 更新时间
draft: false
tags: ["Hugo"]
---

文章内容...
```

每次更新文章时，手动修改 `lastmod` 为当前时间。

### 方法 2：自动使用 Git 提交时间

如果想让 Hugo 自动使用 Git 的最后提交时间作为更新时间，需要：

#### 1. 启用 Git 信息

在 `hugo.toml` 中添加：

```toml
enableGitInfo = true
```

#### 2. 注意事项

- 这会让 Hugo 读取每个文件的 Git 提交历史
- 首次构建会稍慢一些
- 适合团队协作，自动跟踪文件修改时间

### 方法 3：只在文章列表显示更新时间

如果不想改排序，只想在文章页面显示更新时间，只需在 Front Matter 中添加 `lastmod` 即可。

Stack 主题会自动在文章底部显示：

> 📅 最后更新于：2025-02-20 15:30

## 🔍 验证效果

### 查看文章排序

```bash
# 本地预览
hugo serve

# 访问 http://localhost:1313
# 首页的文章列表现在会按更新时间排序
```

### 批量添加 lastmod

如果你想给所有文章批量添加 `lastmod` 字段：

```bash
# 使用脚本批量处理（示例）
for file in content/posts/*.md; do
  # 获取文件的最后修改时间
  lastmod=$(git log -1 --format="%aI" "$file" 2>/dev/null || date -Iseconds)
  
  # 检查是否已有 lastmod
  if ! grep -q "^lastmod:" "$file"; then
    # 在 Front Matter 中插入 lastmod
    sed -i '' "/^date:/a\\
lastmod: $lastmod
" "$file"
  fi
done
```

## 📊 排序优先级总结

现在文章列表的排序规则：

1. ✅ 如果文章有 `lastmod` → 按 `lastmod` 排序
2. ✅ 如果启用了 `enableGitInfo` → 按 Git 提交时间排序
3. ✅ 如果都没有 → 按 `date` 创建时间排序

## 🎯 示例：更新一篇旧文章

假设你有一篇旧文章：

```markdown
---
title: "Hugo 入门教程"
date: 2024-01-15T10:00:00+08:00
draft: false
---
```

现在你更新了内容，添加 `lastmod`：

```markdown
---
title: "Hugo 入门教程"
date: 2024-01-15T10:00:00+08:00
lastmod: 2025-02-20T16:00:00+08:00  # 👈 添加这一行
draft: false
---
```

这篇文章就会出现在首页文章列表的最前面！

## 💡 最佳实践

### 新文章

创建新文章时：

```bash
hugo new posts/my-new-post.md
```

Hugo 会自动设置 `date` 为当前时间。

### 更新文章

每次修改文章内容后：

```bash
# 手动更新 lastmod
# 或者依赖 Git 提交时间（如果启用了 enableGitInfo）
git add content/posts/my-post.md
git commit -m "Update: my post"
```

### RSS Feed

RSS 订阅也会按照更新时间排序，读者能看到你最新更新的文章。

## 🔗 相关链接

- [Hugo Front Matter 文档](https://gohugo.io/content-management/front-matter/)
- [Hugo Git Info 文档](https://gohugo.io/variables/git/)
- [Hugo 日期格式](https://gohugo.io/functions/dateformat/)

