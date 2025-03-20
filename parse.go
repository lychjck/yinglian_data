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

func get_yinglian(text string) {
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
					// _, err := db.Exec(
					// 	"INSERT INTO couplets (first, second, author, dynasty) VALUES (?, ?, ?, ?)",
					// 	upper, lower, "", "",
					// )
					// if err != nil {
					// 	log.Printf("[Error] Failed to insert couplet: %v", err)
					// 	continue
					// }
					// log.Printf("[Success] Successfully inserted couplet: %s, %s", upper, lower)
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		get_yinglian(line)
		fmt.Println("----------------------------------------")
	}

	fmt.Println(count)
}

// 智能校验函数（可根据需要扩展）
func isValidCouplet(upper, lower string) bool {
	// 基础校验
	if len(upper) < 4 || len(lower) < 4 || len(upper) != len(lower){
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
