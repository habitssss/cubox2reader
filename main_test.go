package main

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// 测试HTML解析和收集箱书签提取
func TestFindInboxAndExtractBookmarks(t *testing.T) {
	// 创建一个简单的测试HTML字符串
	htmlStr := `
	<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
	<Title>Bookmarks</Title>
	<H1>Bookmarks</H1>
	<DL><p>
	<DT><H3 FOLDED coverAdaptive="false"  coverType="-1" >收集箱</H3>
	<DL><p>
	<DT><A HREF="https://example.com/test1">测试书签1</A>
	<DD>这是测试书签1的描述
	<DT><A HREF="https://example.com/test2">测试书签2</A>
	<DD>这是测试书签2的描述
	</DL><p>
	</DL><p>
	`

	// 解析HTML
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		t.Fatalf("解析HTML失败: %v", err)
	}

	// 提取书签
	var bookmarks []Bookmark
	var inboxFound bool
	findInboxAndExtractBookmarks(doc, &bookmarks, &inboxFound)

	// 验证结果
	if !inboxFound {
		t.Error("未找到收集箱标签")
	}

	if len(bookmarks) != 2 {
		t.Errorf("期望提取2个书签，但得到了%d个", len(bookmarks))
	}

	// 验证第一个书签
	if len(bookmarks) > 0 {
		firstBookmark := bookmarks[0]
		if firstBookmark.Title != "测试书签1" {
			t.Errorf("第一个书签标题错误，期望 '测试书签1'，得到 '%s'", firstBookmark.Title)
		}
		if firstBookmark.URL != "https://example.com/test1" {
			t.Errorf("第一个书签URL错误，期望 'https://example.com/test1'，得到 '%s'", firstBookmark.URL)
		}
		if firstBookmark.DocumentTags != "这是测试书签1的描述" {
			t.Errorf("第一个书签标签错误，期望 '这是测试书签1的描述'，得到 '%s'", firstBookmark.DocumentTags)
		}
	}

	// 验证第二个书签
	if len(bookmarks) > 1 {
		secondBookmark := bookmarks[1]
		if secondBookmark.Title != "测试书签2" {
			t.Errorf("第二个书签标题错误，期望 '测试书签2'，得到 '%s'", secondBookmark.Title)
		}
		if secondBookmark.URL != "https://example.com/test2" {
			t.Errorf("第二个书签URL错误，期望 'https://example.com/test2'，得到 '%s'", secondBookmark.URL)
		}
		if secondBookmark.DocumentTags != "这是测试书签2的描述" {
			t.Errorf("第二个书签标签错误，期望 '这是测试书签2的描述'，得到 '%s'", secondBookmark.DocumentTags)
		}
	}
}

// 测试指定标签书签提取
func TestFindTagAndExtractBookmarks(t *testing.T) {
	// 创建一个包含多个标签的测试HTML字符串
	htmlStr := `
	<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
	<Title>Bookmarks</Title>
	<H1>Bookmarks</H1>
	<DL><p>
	<DT><H3 FOLDED coverAdaptive="false"  coverType="-1" >收集箱</H3>
	<DL><p>
	<DT><A HREF="https://example.com/inbox1">收集箱书签1</A>
	<DD>收集箱描述1
	</DL><p>
	<DT><H3 FOLDED coverAdaptive="false"  coverType="-1" >技术文章</H3>
	<DL><p>
	<DT><A HREF="https://example.com/tech1">技术文章1</A>
	<DD>技术文章描述1
	<DT><A HREF="https://example.com/tech2">技术文章2</A>
	<DD>技术文章描述2
	</DL><p>
	</DL><p>
	`

	// 解析HTML
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		t.Fatalf("解析HTML失败: %v", err)
	}

	// 提取指定标签的书签
	var bookmarks []Bookmark
	var tagFound bool
	findTagAndExtractBookmarks(doc, &bookmarks, &tagFound, "技术文章")

	// 验证结果
	if !tagFound {
		t.Error("未找到'技术文章'标签")
	}

	if len(bookmarks) != 2 {
		t.Errorf("期望提取2个书签，但得到了%d个", len(bookmarks))
	}

	// 验证书签内容
	if len(bookmarks) > 0 {
		firstBookmark := bookmarks[0]
		if firstBookmark.Title != "技术文章1" {
			t.Errorf("第一个书签标题错误，期望 '技术文章1'，得到 '%s'", firstBookmark.Title)
		}
		if firstBookmark.URL != "https://example.com/tech1" {
			t.Errorf("第一个书签URL错误，期望 'https://example.com/tech1'，得到 '%s'", firstBookmark.URL)
		}
	}
}

// 测试提取所有标签的书签
func TestExtractAllBookmarks(t *testing.T) {
	// 创建一个包含多个标签的测试HTML字符串
	htmlStr := `
	<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
	<Title>Bookmarks</Title>
	<H1>Bookmarks</H1>
	<DL><p>
	<DT><H3 FOLDED coverAdaptive="false"  coverType="-1" >收集箱</H3>
	<DL><p>
	<DT><A HREF="https://example.com/inbox1">收集箱书签1</A>
	<DD>收集箱描述1
	</DL><p>
	<DT><H3 FOLDED coverAdaptive="false"  coverType="-1" >技术文章</H3>
	<DL><p>
	<DT><A HREF="https://example.com/tech1">技术文章1</A>
	<DD>技术文章描述1
	<DT><A HREF="https://example.com/tech2">技术文章2</A>
	<DD>技术文章描述2
	</DL><p>
	</DL><p>
	`

	// 解析HTML
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		t.Fatalf("解析HTML失败: %v", err)
	}

	// 提取所有标签的书签
	var bookmarks []Bookmark
	extractAllBookmarks(doc, &bookmarks)

	// 验证结果
	if len(bookmarks) != 3 {
		t.Errorf("期望提取3个书签，但得到了%d个", len(bookmarks))
	}

	// 验证是否包含了所有标签的书签
	var hasInboxBookmark, hasTechBookmark1, hasTechBookmark2 bool

	for _, bookmark := range bookmarks {
		switch bookmark.URL {
		case "https://example.com/inbox1":
			hasInboxBookmark = true
		case "https://example.com/tech1":
			hasTechBookmark1 = true
		case "https://example.com/tech2":
			hasTechBookmark2 = true
		}
	}

	if !hasInboxBookmark {
		t.Error("未找到收集箱书签")
	}
	if !hasTechBookmark1 {
		t.Error("未找到技术文章1书签")
	}
	if !hasTechBookmark2 {
		t.Error("未找到技术文章2书签")
	}
}

// 测试CSV生成
func TestWriteToCSV(t *testing.T) {
	// 创建测试书签
	bookmarks := []Bookmark{
		{
			Title:           "测试书签1",
			URL:             "https://example.com/test1",
			ID:              "01abcdef1234",
			DocumentTags:    "测试标签1",
			SavedDate:       "2025-01-01 12:00:00.000000+00:00",
			ReadingProgress: "0",
			Location:        "new",
			Seen:            "True",
		},
		{
			Title:           "测试书签2",
			URL:             "https://example.com/test2",
			ID:              "01abcdef5678",
			DocumentTags:    "测试标签2",
			SavedDate:       "2025-01-01 12:00:00.000000+00:00",
			ReadingProgress: "0",
			Location:        "new",
			Seen:            "True",
		},
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	tmpFileName := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpFileName) // 测试结束后删除临时文件

	// 写入CSV
	err = writeToCSV(tmpFileName, bookmarks)
	if err != nil {
		t.Fatalf("写入CSV失败: %v", err)
	}

	// 读取生成的CSV文件
	content, err := os.ReadFile(tmpFileName)
	if err != nil {
		t.Fatalf("读取CSV文件失败: %v", err)
	}

	// 验证CSV内容
	csvStr := string(content)
	if !strings.Contains(csvStr, "Title,URL,ID,Document tags,Saved date,Reading progress,Location,Seen") {
		t.Error("CSV头部不正确")
	}
	if !strings.Contains(csvStr, "测试书签1,https://example.com/test1,01abcdef1234,测试标签1") {
		t.Error("第一个书签数据不正确")
	}
	if !strings.Contains(csvStr, "测试书签2,https://example.com/test2,01abcdef5678,测试标签2") {
		t.Error("第二个书签数据不正确")
	}
}

// 测试ID生成
func TestGenerateID(t *testing.T) {
	id1 := generateID()
	id2 := generateID()

	// 验证ID格式
	if !strings.HasPrefix(id1, "01") {
		t.Errorf("生成的ID格式不正确，应以'01'开头，得到: %s", id1)
	}

	// 验证ID的唯一性
	if id1 == id2 {
		t.Error("生成的两个ID相同，应该是唯一的")
	}
}