package main

import (

	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"flag"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func parse_yinglian() {
	var filename string
	flag.StringVar(&filename,"name","","name")
	flag.Parse()

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("文件读取失败:", err)
	}

	dbName := "couples"
	dbHost := "localhost"
	dbUser := "root"
	dbPass := ""
	log.Printf("[Database] Connecting to MySQL database %s on %s", dbName, dbHost)
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName+"?parseTime=true&loc=Local&charset=utf8mb4")
	if err != nil {
		log.Fatal("[Error] Database connection failed: ", err)
	}
	log.Println("[Database] Connected successfully")
	defer db.Close()

	text := string(content)

	// 复杂匹配模式：支持三种引号格式和不同分隔符
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

	// 存储已发现的对联避免重复
	uniqueSet := make(map[string]bool)

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}

			// 清洗数据并格式化
			upper := strings.TrimSpace(match[1])
			lower := strings.TrimSpace(match[2])

			// 智能过滤条件
			if isValidCouplet(upper, lower) {
				key := upper + "|" + lower
				if !uniqueSet[key] {
					// fmt.Printf("上联：%s\n下联：%s\n\n", upper, lower)
					uniqueSet[key] = true
					_, err := db.Exec(
						"INSERT INTO couplets (first, second, author, dynasty) VALUES (?, ?, ?, ?)",
						upper, lower, "", "",
					)
					if err != nil {
						log.Printf("[Error] Failed to insert couplet: %v", err)
						continue
					}
					log.Printf("[Success] Successfully inserted couplet: %s, %s", upper, lower)
				}
			}
		}
	}
}

// 智能校验函数（可根据需要扩展）
func isValidCouplet(upper, lower string) bool {
	// 基础校验
	if len(upper) < 4 || len(lower) < 4 {
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