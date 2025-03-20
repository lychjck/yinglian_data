package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// 存储已发现的对联避免重复
var uniqueSet = make(map[string]bool)
var count = 0

func get_yinglian(text string, db *sql.DB, ref int64,source string) {
	patterns := []*regexp.Regexp{
		// 匹配全角引号对：上下联完整包含在引号内
		regexp.MustCompile(`【([^】]+)[;；]([^】]+)】`),
		// 匹配直角引号对
		regexp.MustCompile(`「([^」]+)[;；]([^」]+)」`),
		// 匹配传统中文引号对
		regexp.MustCompile(`“([^”]+)[;；]([^”]+)”`),
		// 匹配无引号但对仗工整的句式
		regexp.MustCompile(`([\p{Han}，、]{5,}?)[;；]([\p{Han}，、]{5,})`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}

			// 清洗数据并格式化
			upper := strings.TrimSpace(match[1])
			lower := strings.TrimSpace(match[2])
			lower = strings.ReplaceAll(lower, "。", "")
			lower = strings.ReplaceAll(lower, "？", "")
			lower = strings.ReplaceAll(lower, "！", "")

			// 智能过滤条件
			if isValidCouplet(upper, lower) {
				key := upper + "|" + lower
				if !uniqueSet[key] {
					fmt.Printf("上联：%s\n下联：%s\n\n", upper, lower)
					uniqueSet[key] = true
					count++
					_, err := db.Exec(
						"INSERT INTO yinglian (first, second, ref,source) VALUES (?, ?, ?,?)",
						upper, lower, ref, source,
					)
					if err != nil {
						log.Printf("[Error] Failed to insert yinglian: %v", err)
						continue
					}
				}
			}
		}
	}
}

func parse_yinglian(filename string, db *sql.DB) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("文件打开失败:", err)
	}
	defer file.Close()

	var bookName, volume, title string
	lineCount := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// 处理前三行特殊信息
		switch lineCount {
		case 1:
			bookName = line
		case 2:
			volume = line
		case 3:
			title = line
		default:

			if bookName == "楹联三话" || bookName == "散见联话"{
				if lineCount % 2 == 0 {
					title = strings.TrimSpace(line)
					title = strings.ReplaceAll(title, " ", "")
					title = strings.ReplaceAll(title, "\t", "")
					title = strings.ReplaceAll(title, "  ", "")
					continue
				}
			}
			// 处理正文内容
			fmt.Println(line)
			
			// 将内容插入到yinglian_content表
			result, err := db.Exec(
				"INSERT INTO yinglian_content (book_name, volume, title, content) VALUES (?, ?, ?, ?)",
				bookName, volume, title, line,
			)
			if err != nil {
				log.Printf("[Error] Failed to insert content: %v", err)
				continue
			}

			// 获取插入的ID
			contentId, err := result.LastInsertId()
			if err != nil {
				log.Printf("[Error] Failed to get last insert id: %v", err)
				continue
			}
			get_yinglian(line, db, contentId,bookName)
			fmt.Println("----------------------------------------")
		}
	}

	fmt.Printf("处理完成，共发现%d个对联\n", count)
}

// 智能校验函数（可根据需要扩展）
func isValidCouplet(upper, lower string) bool {
	// 基础校验
	if len(upper) < 4 || len(lower) < 4 || len(upper) != len(lower) {
		return false
	}

	// 对仗特征校验
	features := []struct {
		check func(string) bool
		count int
	}{
		{check: func(s string) bool { return strings.Contains(s, "，") }, count: 2},
		{check: func(s string) bool { return strings.ContainsAny(s, "仄平") }, count: 1},
	}

	score := 0
	for _, feat := range features {
		if feat.check(upper) == feat.check(lower) {
			score += feat.count
		}
	}
	return score >= 2
}
