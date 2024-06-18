package gospider

import (
	"bytes"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// NewXpathParser Xpath构造函数
func NewXpathParser(content []byte) Xpath {
	doc, _ := htmlquery.Parse(bytes.NewReader(content))
	return Xpath{doc}
}

// Xpath Xpath解析html
type Xpath struct {
	X *html.Node
}

// XpathList 获取所有符合条件的节点
func (x *Xpath) XpathList(s string) []Xpath {
	nodes := htmlquery.Find(x.X, s)
	result := []Xpath{}
	if x.X == nil {
		return result
	}
	for _, node := range nodes {
		result = append(result, Xpath{X: node})
	}
	return result
}

// Xpath 获取符合条件的第一个节点
func (x *Xpath) Xpath(s string) Xpath {
	node := htmlquery.FindOne(x.X, s)
	return Xpath{X: node}
}

// Extract 获取符合条件的所有节点的文本内容
func (x *Xpath) Extract(s string) []string {
	result := []string{}
	if x.X == nil {
		return result
	}
	nodes := htmlquery.Find(x.X, s)
	for _, node := range nodes {
		result = append(result, htmlquery.InnerText(node))
	}
	return result
}

// ExtractFirst 获取符合条件的第一个节点的文本内容
func (x *Xpath) ExtractFirst(s string) string {
	var ef string
	if x.X == nil {
		return ef
	}
	node := htmlquery.FindOne(x.X, s)
	if node != nil {
		return htmlquery.InnerText(node)
	} else {
		return ef
	}
}
