# Cubox转Reader

这个工具用于将Cubox导出的书签HTML文件转换为Reader可导入的CSV格式。

## 功能特点

- 解析Cubox导出的HTML书签文件
- 可以提取所有标签下的书签内容，或指定特定标签
- 生成符合Reader导入要求的CSV格式文件
- 自动生成必要的ID和时间戳

## 使用方法

### 安装

```bash
# 克隆仓库
git clone https://github.com/habitssss/cubox2reader.git
cd cubox2reader

# 构建项目
go build
```

### 命令行参数

```
用法: cubox2reader -html=<html文件路径> [-output=<输出csv文件>] [-tag=<标签名称>]
  -html string
        HTML文件路径 (必填)
  -output string
        输出CSV文件路径 (默认 "output.csv")
  -tag string
        指定要处理的标签名称 (不指定则处理所有标签)
```

### 示例

```bash
# 基本使用 (处理所有标签)
./cubox2reader -html=cubox_bookmarks.html

# 指定输出文件
./cubox2reader -html=cubox_bookmarks.html -output=reader_import.csv

# 只处理特定标签
./cubox2reader -html=cubox_bookmarks.html -tag="收集箱"

# 处理特定标签并指定输出文件
./cubox2reader -html=cubox_bookmarks.html -tag="技术文章" -output=tech_articles.csv
```

## CSV输出格式

生成的CSV文件包含以下字段：

- Title: 书签标题
- URL: 书签链接地址
- ID: 唯一标识符
- Document tags: 文档标签
- Saved date: 保存日期（UTC格式）
- Reading progress: 阅读进度（默认为0）
- Location: 位置（默认为"new"）
- Seen: 是否查看过（默认为"False"）

## 注意事项

- HTML文件必须是Cubox导出的标准书签格式
- 默认情况下程序会提取所有标签下的内容，可以通过 `-tag` 参数指定只处理特定标签
- 生成的ID采用基于时间戳的格式，确保唯一性
- 如果指定的标签不存在，程序会给出警告信息

## 工作原理

1. 解析Cubox导出的HTML书签文件
2. 根据用户指定的参数，提取所有标签或特定标签下的书签
3. 为每个书签生成唯一ID和其他必要字段
4. 将处理后的数据以CSV格式保存，符合Reader的导入要求

## 许可证

MIT