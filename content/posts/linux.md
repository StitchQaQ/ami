---
title: "Linux 文本处理四大神器：grep、awk、sed、jq"
date: 2025-10-30T16:55:54+08:00
lastmod: 2025-10-30T17:30:00+08:00
draft: false
description: "深入理解 Linux 文本处理工具 grep、awk、sed、jq 的实用技巧"
tags: ["Linux", "命令行", "文本处理"]
categories: ["工具"]
author: "Heyuuuu"
toc: true
---

## 📚 概述

在 Linux 系统中，文本处理是日常运维和开发的核心技能。掌握这四大神器，能让你的工作效率提升 10 倍！

- **grep**："搜索高手"，用于查找文本中的匹配模式
- **awk**："数据专家"，用于从结构化文本中提取信息并计算
- **sed**："流编辑器"，用于批量修改文本内容
- **jq**："JSON 利器"，用于过滤和转换 JSON 数据

---

## 1. grep - 文本搜索神器

### 🎯 核心功能

`grep` (Global Regular Expression Print) 用于在文件中搜索匹配指定模式的行。

### 📖 基础用法

```bash
# 基本语法
grep [选项] 模式 文件名
```

### 💡 常用选项

| 选项 | 说明 | 示例 |
|------|------|------|
| `-i` | 忽略大小写 | `grep -i "error" log.txt` |
| `-n` | 显示行号 | `grep -n "ERROR" log.txt` |
| `-v` | 反向匹配（不包含） | `grep -v "DEBUG" log.txt` |
| `-r` | 递归搜索目录 | `grep -r "TODO" ./src` |
| `-c` | 统计匹配行数 | `grep -c "ERROR" log.txt` |
| `-l` | 只显示文件名 | `grep -l "config" *.conf` |
| `-A n` | 显示匹配行及后n行 | `grep -A 3 "ERROR" log.txt` |
| `-B n` | 显示匹配行及前n行 | `grep -B 2 "ERROR" log.txt` |
| `-C n` | 显示匹配行及前后n行 | `grep -C 2 "ERROR" log.txt` |
| `-E` | 使用扩展正则表达式 | `grep -E "error|warning" log.txt` |

### 🔥 实战案例

#### 案例 1：日志分析

```bash
# 查找所有错误日志
grep "ERROR" server.log

# 查找错误和警告（不区分大小写）
grep -i -E "error|warning" server.log

# 统计今天的错误数量
grep "2025-10-30.*ERROR" server.log | wc -l

# 查找错误并显示前后3行上下文
grep -C 3 "ERROR" server.log

# 查找包含特定IP的访问记录
grep "192.168.1.100" access.log
```

#### 案例 2：代码搜索

```bash
# 在所有Python文件中查找TODO注释
grep -rn "TODO" --include="*.py" .

# 查找函数定义
grep -rn "def.*login" --include="*.py" .

# 排除某些目录
grep -r "import" --exclude-dir={venv,node_modules} .

# 查找空行
grep -n "^$" config.py

# 查找以特定字符开头的行
grep "^class " models.py
```

#### 案例 3：系统运维

```bash
# 查找正在运行的进程
ps aux | grep python

# 查找某个端口被哪个进程占用
netstat -tulnp | grep :8080

# 查看活跃的SSH连接
who | grep pts

# 查找大文件
ls -lh | grep "G\|M"
```

### 🚀 高级技巧

```bash
# 使用正则表达式匹配邮箱
grep -E "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}" contacts.txt

# 匹配IP地址
grep -E "\b([0-9]{1,3}\.){3}[0-9]{1,3}\b" network.log

# 多个模式匹配（OR）
grep -E "error|fail|exception" app.log

# 排除多个模式
grep -v -E "DEBUG|INFO" app.log

# 统计每个错误类型的数量
grep -o "ERROR: [^:]*" log.txt | sort | uniq -c
```

---

## 2. awk - 数据处理专家

### 🎯 核心功能

`awk` 是一个强大的文本分析工具，特别擅长处理格式化的文本数据（如CSV、日志等）。

### 📖 基础语法

```bash
awk 'pattern { action }' 文件名
```

### 💡 内置变量

| 变量 | 说明 |
|------|------|
| `$0` | 当前整行 |
| `$1, $2, $3...` | 第1列、第2列、第3列... |
| `NR` | 当前行号 |
| `NF` | 当前行的字段数 |
| `FS` | 字段分隔符（默认空格） |
| `OFS` | 输出字段分隔符 |

### 🔥 实战案例

#### 案例 1：日志分析

```bash
# 打印日志的第1列和第5列
awk '{print $1, $5}' access.log

# 打印行号和内容
awk '{print NR, $0}' file.txt

# 统计访问日志中的状态码分布
awk '{print $9}' access.log | sort | uniq -c | sort -rn

# 计算响应时间的平均值（假设在第10列）
awk '{sum+=$10; count++} END {print "平均响应时间:", sum/count}' access.log

# 筛选响应时间大于1秒的请求
awk '$10 > 1 {print $0}' access.log
```

#### 案例 2：CSV 数据处理

假设有一个 `employees.csv` 文件：

```csv
name,age,salary,department
Alice,28,5000,IT
Bob,35,6000,HR
Charlie,32,5500,IT
David,29,4800,Finance
```

```bash
# 使用逗号作为分隔符
awk -F',' '{print $1, $3}' employees.csv

# 跳过表头，只处理数据行
awk -F',' 'NR>1 {print $1, $3}' employees.csv

# 计算IT部门的平均工资
awk -F',' 'NR>1 && $4=="IT" {sum+=$3; count++} END {print "IT平均工资:", sum/count}' employees.csv

# 筛选工资大于5000的员工
awk -F',' 'NR>1 && $3>5000 {print $1, $3}' employees.csv

# 格式化输出
awk -F',' 'NR>1 {printf "%-10s %5d %10s\n", $1, $2, $4}' employees.csv
```

#### 案例 3：系统监控

```bash
# 显示CPU使用率最高的10个进程
ps aux | awk 'NR>1 {print $11, $3}' | sort -k2 -rn | head -10

# 统计内存使用情况
free -m | awk 'NR==2 {printf "内存使用率: %.2f%%\n", $3/$2*100}'

# 分析磁盘使用情况
df -h | awk 'NR>1 {if($5+0 > 80) print $1, $5, "磁盘空间不足!"}'

# 统计网络连接状态
netstat -an | awk '/^tcp/ {state[$6]++} END {for(i in state) print i, state[i]}'
```

### 🚀 高级技巧

```bash
# 使用条件判断
awk '{if($3 > 5000) print $1, "高薪"; else print $1, "普通"}' employees.csv

# 使用循环
awk '{for(i=1; i<=NF; i++) print $i}' file.txt

# 使用数组统计
awk '{count[$1]++} END {for(ip in count) print ip, count[ip]}' access.log

# 多分隔符
awk -F'[,:]' '{print $1, $3}' data.txt

# 自定义输出分隔符
awk 'BEGIN{OFS="|"} {print $1, $2, $3}' file.txt

# 处理多个文件
awk '{print FILENAME, $0}' file1.txt file2.txt
```

---

## 3. sed - 流编辑器

### 🎯 核心功能

`sed` (Stream Editor) 用于对文本进行批量替换、删除、插入等操作。

### 📖 基础语法

```bash
sed [选项] '命令' 文件名
```

### 💡 常用命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `s/old/new/` | 替换 | `sed 's/foo/bar/' file.txt` |
| `s/old/new/g` | 全局替换 | `sed 's/foo/bar/g' file.txt` |
| `d` | 删除行 | `sed '3d' file.txt` |
| `p` | 打印行 | `sed -n '1,5p' file.txt` |
| `a` | 追加行 | `sed '3a新行内容' file.txt` |
| `i` | 插入行 | `sed '3i新行内容' file.txt` |
| `c` | 替换整行 | `sed '3c新内容' file.txt` |

### 🔥 实战案例

#### 案例 1：文本替换

```bash
# 替换第一次出现的模式
sed 's/old/new/' file.txt

# 全局替换（每行所有匹配）
sed 's/old/new/g' file.txt

# 替换并保存到原文件
sed -i 's/old/new/g' file.txt

# 替换指定行
sed '3s/old/new/' file.txt

# 替换范围内的行
sed '1,10s/old/new/g' file.txt

# 大小写不敏感替换
sed 's/error/ERROR/gi' log.txt
```

#### 案例 2：删除操作

```bash
# 删除第3行
sed '3d' file.txt

# 删除最后一行
sed '$d' file.txt

# 删除1到5行
sed '1,5d' file.txt

# 删除空行
sed '/^$/d' file.txt

# 删除包含特定文本的行
sed '/DEBUG/d' log.txt

# 删除以#开头的注释行
sed '/^#/d' config.txt
```

#### 案例 3：插入和追加

```bash
# 在第3行后追加内容
sed '3a这是新增的行' file.txt

# 在第3行前插入内容
sed '3i这是插入的行' file.txt

# 在匹配行后追加
sed '/pattern/a新行内容' file.txt

# 在文件开头插入
sed '1i文件头部' file.txt

# 在文件末尾追加
sed '$a文件尾部' file.txt
```

#### 案例 4：配置文件修改

```bash
# 修改配置文件中的端口
sed -i 's/port=8080/port=9000/g' config.ini

# 注释掉某一行
sed -i '5s/^/#/' config.txt

# 取消注释
sed -i 's/^#//' config.txt

# 替换IP地址
sed -i 's/192\.168\.1\.100/192.168.1.200/g' network.conf

# 批量修改多个文件
sed -i 's/old/new/g' *.conf
```

### 🚀 高级技巧

```bash
# 使用正则表达式捕获组
sed 's/\(.*\)@\(.*\)/用户:\1 域名:\2/' emails.txt

# 多个替换操作
sed -e 's/foo/bar/g' -e 's/hello/hi/g' file.txt

# 使用不同的分隔符
sed 's|/usr/local|/opt|g' paths.txt

# 只打印匹配的行（-n 参数）
sed -n '/ERROR/p' log.txt

# 替换特定行号范围
sed '10,20s/old/new/g' file.txt

# 条件替换
sed '/pattern/s/old/new/g' file.txt

# 删除行尾空格
sed 's/[[:space:]]*$//' file.txt

# 每行前添加行号
sed = file.txt | sed 'N;s/\n/\t/'
```

---

## 4. jq - JSON 处理利器

### 🎯 核心功能

`jq` 是专门用于处理 JSON 数据的命令行工具，类似于 JSON 的 `grep`。

### 📖 基础语法

```bash
jq [选项] '过滤器' 文件名
```

### 💡 常用过滤器

| 过滤器 | 说明 | 示例 |
|--------|------|------|
| `.` | 整个JSON | `jq '.' data.json` |
| `.key` | 获取键值 | `jq '.name' data.json` |
| `.[]` | 数组元素 | `jq '.users[]' data.json` |
| `.[n]` | 第n个元素 | `jq '.[0]' array.json` |
| `.key.subkey` | 嵌套访问 | `jq '.user.name' data.json` |

### 🔥 实战案例

#### 案例 1：基本查询

假设有一个 `users.json` 文件：

```json
{
  "users": [
    {"id": 1, "name": "Alice", "age": 28, "city": "Beijing"},
    {"id": 2, "name": "Bob", "age": 35, "city": "Shanghai"},
    {"id": 3, "name": "Charlie", "age": 32, "city": "Beijing"}
  ]
}
```

```bash
# 美化输出整个JSON
jq '.' users.json

# 获取所有用户
jq '.users' users.json

# 获取第一个用户
jq '.users[0]' users.json

# 获取所有用户的名字
jq '.users[].name' users.json

# 获取特定字段
jq '.users[] | {name, age}' users.json
```

#### 案例 2：过滤和筛选

```bash
# 筛选年龄大于30的用户
jq '.users[] | select(.age > 30)' users.json

# 筛选来自北京的用户
jq '.users[] | select(.city == "Beijing")' users.json

# 筛选并只显示名字
jq '.users[] | select(.age > 30) | .name' users.json

# 多条件筛选
jq '.users[] | select(.age > 25 and .city == "Beijing")' users.json

# 排除某些条件
jq '.users[] | select(.city != "Shanghai")' users.json
```

#### 案例 3：数组操作

```bash
# 获取数组长度
jq '.users | length' users.json

# 获取前2个元素
jq '.users[:2]' users.json

# 排序（按年龄）
jq '.users | sort_by(.age)' users.json

# 反向排序
jq '.users | sort_by(.age) | reverse' users.json

# 去重
jq '.users | unique_by(.city)' users.json

# 分组统计
jq '.users | group_by(.city) | map({city: .[0].city, count: length})' users.json
```

#### 案例 4：API 响应处理

```bash
# 从 API 获取数据并提取字段
curl -s 'https://api.github.com/users/github' | jq '.name, .public_repos'

# 提取数组中的特定字段
curl -s 'https://api.github.com/users/github/repos' | jq '.[].name'

# 统计 star 数最多的仓库
curl -s 'https://api.github.com/users/github/repos' | jq 'sort_by(.stargazers_count) | reverse | .[0]'

# 格式化输出
curl -s 'https://api.github.com/users/github/repos' | jq '.[] | "\(.name): \(.stargazers_count) stars"'
```

#### 案例 5：数据转换

```bash
# 构造新的JSON对象
jq '.users[] | {姓名: .name, 年龄: .age}' users.json

# 添加新字段
jq '.users[] | . + {status: "active"}' users.json

# 修改字段值
jq '.users[] | .age = .age + 1' users.json

# 重命名字段
jq '.users[] | {id, username: .name, years: .age}' users.json

# 合并对象
jq '. + {timestamp: "2025-10-30"}' data.json
```

### 🚀 高级技巧

```bash
# 条件判断
jq '.users[] | if .age > 30 then "老员工" else "新员工" end' users.json

# 使用变量
jq --arg city "Beijing" '.users[] | select(.city == $city)' users.json

# 读取多个JSON文件
jq -s '.[0].users + .[1].users' file1.json file2.json

# 输出为CSV格式
jq -r '.users[] | [.name, .age, .city] | @csv' users.json

# 从CSV转JSON（需要配合其他工具）
jq -R -s 'split("\n") | map(split(","))' data.csv

# 递归搜索
jq '.. | select(type == "string" and contains("Beijing"))' users.json

# 错误处理
jq '.users[]? // "默认值"' users.json

# 复杂计算
jq '.users | map(.age) | add / length' users.json  # 计算平均年龄
```

---

## 5. 组合使用技巧

### 🔗 管道组合

```bash
# grep + awk: 提取错误日志的时间戳
grep "ERROR" log.txt | awk '{print $1, $2}'

# grep + sed: 查找并替换
grep -l "old_text" *.txt | xargs sed -i 's/old_text/new_text/g'

# awk + sort + uniq: 统计访问最多的IP
awk '{print $1}' access.log | sort | uniq -c | sort -rn | head -10

# curl + jq: API数据提取
curl -s https://api.example.com/data | jq '.results[] | select(.status == "active")'

# find + grep: 在项目中搜索特定代码
find . -name "*.py" -exec grep -l "def login" {} \;
```

### 📝 实用脚本示例

#### 日志分析脚本

```bash
#!/bin/bash
# 分析Nginx访问日志

LOG_FILE="/var/log/nginx/access.log"

echo "=== Nginx日志分析 ==="
echo

echo "1. 访问最多的10个IP:"
awk '{print $1}' $LOG_FILE | sort | uniq -c | sort -rn | head -10

echo
echo "2. 访问最多的10个URL:"
awk '{print $7}' $LOG_FILE | sort | uniq -c | sort -rn | head -10

echo
echo "3. HTTP状态码分布:"
awk '{print $9}' $LOG_FILE | sort | uniq -c | sort -rn

echo
echo "4. 平均响应时间:"
awk '{sum+=$NF; count++} END {print sum/count "秒"}' $LOG_FILE
```

#### JSON数据处理脚本

```bash
#!/bin/bash
# 处理GitHub API返回的仓库数据

USERNAME="github"

echo "获取 $USERNAME 的仓库信息..."

curl -s "https://api.github.com/users/$USERNAME/repos" | \
jq -r '.[] | "\(.name) - \(.stargazers_count) stars - \(.language // "Unknown")"' | \
column -t -s'-'
```

---

## 6. 常见问题和技巧

### ❓ Q&A

**Q: grep、awk、sed 的区别？**

- **grep**: 只是搜索，不修改文件
- **awk**: 主要用于数据提取和统计
- **sed**: 主要用于文本替换和编辑

**Q: 什么时候用 awk，什么时候用 sed？**

- 需要按列处理数据 → 用 `awk`
- 需要替换文本内容 → 用 `sed`
- 需要复杂计算 → 用 `awk`
- 简单的行删除/插入 → 用 `sed`

**Q: jq 必须安装吗？**

`jq` 需要单独安装：

```bash
# macOS
brew install jq

# Ubuntu/Debian
sudo apt install jq

# CentOS/RHEL
sudo yum install jq
```

### 💡 性能优化技巧

1. **使用管道避免中间文件**

```bash
# ❌ 低效
grep "ERROR" log.txt > temp.txt
awk '{print $1}' temp.txt
rm temp.txt

# ✅ 高效
grep "ERROR" log.txt | awk '{print $1}'
```

2. **合并多个 sed 操作**

```bash
# ❌ 多次读取文件
sed 's/foo/bar/g' file.txt | sed 's/old/new/g' | sed 's/a/b/g'

# ✅ 一次完成
sed -e 's/foo/bar/g' -e 's/old/new/g' -e 's/a/b/g' file.txt
```

3. **使用 awk 替代多个管道**

```bash
# ❌ 多个命令
cat file.txt | grep pattern | cut -d' ' -f1 | sort | uniq -c

# ✅ 用 awk 一步完成
awk '/pattern/ {count[$1]++} END {for(i in count) print count[i], i}' file.txt
```

---

## 7. 学习资源

### 📚 推荐阅读

- **grep**: `man grep` 或 `grep --help`
- **awk**: 《The AWK Programming Language》
- **sed**: 《sed & awk》O'Reilly
- **jq**: [官方文档](https://jqlang.github.io/jq/)

### 🔗 在线工具

- [Regex101](https://regex101.com/) - 正则表达式测试
- [jq play](https://jqplay.org/) - jq 在线测试
- [ExplainShell](https://explainshell.com/) - 命令解释

### 💻 练习建议

1. 每天处理实际的日志文件
2. 尝试用一条命令完成复杂任务
3. 编写自动化脚本
4. 参与开源项目，阅读他人的脚本

---

## 📝 总结

掌握这四大工具，你就拥有了强大的文本处理能力：

- 🔍 **grep** - 快速定位问题
- 📊 **awk** - 深入分析数据
- ✏️ **sed** - 批量修改配置
- 🎯 **jq** - 优雅处理JSON

记住：**实践是最好的老师！** 从今天开始，用这些工具处理你的日常任务吧！
