package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	// 获取所有htm文件
	files, err := filepath.Glob("*.htm")
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		fmt.Println("未找到htm文件")
		return
	}

	for _, f := range files {
		// 读取HTML文件并转换编码
		file, err := os.Open(f)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// 使用GBK解码器读取文件
		reader := transform.NewReader(file, simplifiedchinese.GBK.NewDecoder())
		content, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}

		// 解析HTML
		doc, err := html.Parse(strings.NewReader(string(content)))
		if err != nil {
			panic(err)
		}

		// 创建输出文件
		outName := strings.TrimSuffix(f, path.Ext(f)) + "_clean.txt"
		outFile, err := os.Create(outName)
		if err != nil {
			panic(err)
		}
		defer outFile.Close()

		// 直接使用UTF-8编码写入文件
		writer := outFile

		// 遍历DOM树并提取文本
		var traverse func(*html.Node, bool)
		traverse = func(n *html.Node, skip bool) {
			if n.Type == html.ElementNode && n.Data == "em" {
				skip = true
			}

			if n.Type == html.TextNode && !skip {
				text := strings.TrimSpace(n.Data)
				if text != "" {
					_, err := writer.Write([]byte(text))
					if err != nil {
						panic(err)
					}
					// 检查下一个兄弟节点，如果是换行相关的节点，则添加换行
					if n.NextSibling == nil || (n.NextSibling.Type == html.ElementNode && (n.NextSibling.Data == "br" || n.NextSibling.Data == "p" || n.NextSibling.Data == "div")) {
						_, err = writer.Write([]byte("\n"))
						if err != nil {
							panic(err)
						}
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				traverse(c, skip)
			}
		}

		traverse(doc, false)
		writer.Close()
		fmt.Printf("已处理文件：%s → %s\n", f, outName)
	}

}
