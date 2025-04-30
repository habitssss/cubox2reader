package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// 书签项结构
type Bookmark struct {
	Title            string
	URL              string
	ID               string
	DocumentTags     string
	SavedDate        string
	ReadingProgress  string
	Location         string
	Seen             string
}

func main() {
	// 定义命令行参数
	htmlFile := flag.String("html", "", "HTML文件路径 (必填)")
	outputFile := flag.String("output", "output.csv", "输出CSV文件路径")
	tag := flag.String("tag", "", "指定要处理的标签名称 (不指定则处理所有标签)")
	flag.Parse()

	if *htmlFile == "" {
		fmt.Println("错误: 必须指定HTML文件路径")
		fmt.Println("用法: cubox2reader -html=<html文件路径> [-output=<输出csv文件>] [-tag=<标签名称>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 打开HTML文件
	file, err := os.Open(*htmlFile)
	if err != nil {
		log.Fatalf("无法打开HTML文件: %v", err)
	}
	defer file.Close()

	// 解析HTML
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatalf("无法解析HTML: %v", err)
	}

	// 提取书签
	var bookmarks []Bookmark
	var tagFound bool

	if *tag == "" {
		// 如果未指定标签，处理所有标签
		extractAllBookmarks(doc, &bookmarks)
		tagFound = len(bookmarks) > 0
	} else {
		// 处理指定的标签
		findTagAndExtractBookmarks(doc, &bookmarks, &tagFound, *tag)
	}

	if !tagFound {
		if *tag == "" {
			log.Println("警告: 未找到任何书签")
		} else {
			log.Printf("警告: 未找到标签 '%s'\n", *tag)
		}
	}

	// 写入CSV
	err = writeToCSV(*outputFile, bookmarks)
	if err != nil {
		log.Fatalf("写入CSV失败: %v", err)
	}

	fmt.Printf("成功将%d个书签转换为CSV文件: %s\n", len(bookmarks), *outputFile)
}

// 递归查找指定标签并提取内容
func findTagAndExtractBookmarks(n *html.Node, bookmarks *[]Bookmark, tagFound *bool, tagName string) {
	if n == nil {
		return
	}

	// 查找H3标签，判断是否为指定标签
	if n.Type == html.ElementNode && n.Data == "h3" {
		var isTargetTag bool

		// 检查属性
		for _, a := range n.Attr {
			if a.Key == "folded" {
				// 获取H3标签的文本内容
				var text string
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						text += c.Data
					}
				}

				// 如果是指定标签
				if strings.TrimSpace(text) == tagName {
					isTargetTag = true
					*tagFound = true
					break
				}
			}
		}

		// 如果找到指定标签，提取其下的DL标签内容
		if isTargetTag {
			// 查找父节点的下一个DL节点
			parent := n.Parent
			if parent != nil {
				for sibling := n.NextSibling; sibling != nil; sibling = sibling.NextSibling {
					if sibling.Type == html.ElementNode && sibling.Data == "dl" {
						extractBookmarksFromDL(sibling, bookmarks)
						return
					}
				}
			}
		}
	}

	// 如果没找到指定标签，递归处理子节点
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findTagAndExtractBookmarks(c, bookmarks, tagFound, tagName)
	}
}

// 递归查找收集箱标签并提取内容 (保留向后兼容性)
func findInboxAndExtractBookmarks(n *html.Node, bookmarks *[]Bookmark, inboxFound *bool) {
	var tagFound bool
	findTagAndExtractBookmarks(n, bookmarks, &tagFound, "收集箱")
	*inboxFound = tagFound
}

// 提取所有标签下的书签
func extractAllBookmarks(n *html.Node, bookmarks *[]Bookmark) {
	if n == nil {
		return
	}

	// 查找所有H3标签
	if n.Type == html.ElementNode && n.Data == "h3" {
		// 检查是否有folded属性，表示这是一个标签
		var hasFolder bool
		for _, a := range n.Attr {
			if a.Key == "folded" {
				hasFolder = true
				break
			}
		}

		if hasFolder {
			// 查找父节点的下一个DL节点
			parent := n.Parent
			if parent != nil {
				for sibling := n.NextSibling; sibling != nil; sibling = sibling.NextSibling {
					if sibling.Type == html.ElementNode && sibling.Data == "dl" {
						extractBookmarksFromDL(sibling, bookmarks)
						break
					}
				}
			}
		}
	}

	// 递归处理子节点
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractAllBookmarks(c, bookmarks)
	}
}

// 从DL标签提取书签
func extractBookmarksFromDL(dl *html.Node, bookmarks *[]Bookmark) {
	if dl == nil || dl.Type != html.ElementNode || dl.Data != "dl" {
		return
	}

	var currentTitle, currentURL string
	var descriptionPending bool

	// 遍历DL的子节点
	for node := dl.FirstChild; node != nil; node = node.NextSibling {
		if node.Type == html.ElementNode {
			// 处理DT标签
			if node.Data == "dt" {
				// 查找A标签
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && c.Data == "a" {
						// 获取URL
						for _, attr := range c.Attr {
							if attr.Key == "href" {
								currentURL = attr.Val
								break
							}
						}

						// 获取标题
						var titleText string
						for tc := c.FirstChild; tc != nil; tc = tc.NextSibling {
							if tc.Type == html.TextNode {
								titleText += tc.Data
							}
						}
						currentTitle = titleText
						descriptionPending = true
					}
				}
			} else if node.Data == "dd" && descriptionPending {
				// 从DD标签获取描述
				var description string
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						description += c.Data
					}
				}

				// 创建书签并添加到列表
				if currentTitle != "" && currentURL != "" {
					bookmark := Bookmark{
						Title:           currentTitle,
						URL:             currentURL,
						ID:              generateID(),
						DocumentTags:    strings.TrimSpace(description), // 可以将描述作为标签
						SavedDate:       time.Now().UTC().Format("2006-01-02 15:04:05.000000+00:00"),
						ReadingProgress: "0",
						Location:        "new",
						Seen:            "False",
					}
					*bookmarks = append(*bookmarks, bookmark)
				}

				// 重置
				currentTitle = ""
				currentURL = ""
				descriptionPending = false
			}
		}
	}

	// 处理嵌套的文件夹
	for node := dl.FirstChild; node != nil; node = node.NextSibling {
		if node.Type == html.ElementNode && node.Data == "dl" {
			extractBookmarksFromDL(node, bookmarks)
		}
	}
}

// 用于确保ID唯一性的计数器
var idCounter int64 = 0

// 生成唯一ID
func generateID() string {
	now := time.Now().UTC()
	timestamp := now.UnixNano() / 1000000 // 毫秒级时间戳

	// 增加计数器以确保唯一性
	idCounter++

	// 将时间戳和计数器组合以确保唯一性
	return fmt.Sprintf("01%s%x", strings.ToLower(fmt.Sprintf("%x", timestamp)), idCounter)
}

// 写入CSV
func writeToCSV(outputPath string, bookmarks []Bookmark) error {
	// 创建输出目录
	dir := filepath.Dir(outputPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建输出目录失败: %v", err)
		}
	}

	// 创建CSV文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("无法创建CSV文件: %v", err)
	}
	defer file.Close()

	// 创建CSV写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入CSV头
	headers := []string{
		"Title", "URL", "ID", "Document tags", "Saved date", "Reading progress", "Location", "Seen",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("写入CSV头失败: %v", err)
	}

	// 写入书签数据
	for _, bookmark := range bookmarks {
		record := []string{
			bookmark.Title,
			bookmark.URL,
			bookmark.ID,
			bookmark.DocumentTags,
			bookmark.SavedDate,
			bookmark.ReadingProgress,
			bookmark.Location,
			bookmark.Seen,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("写入CSV记录失败: %v", err)
		}
	}

	return nil
}