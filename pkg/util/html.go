package util

import (
	"golang.org/x/net/html"
	"regexp"
)

func FindTagByNameAttrs(node *html.Node, name string, useAttr bool, attrs []html.Attribute) []*html.Node {
	res := make([]*html.Node, 0)
	q := []*html.Node{node}

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		if cur.Type == html.ElementNode && cur.Data == name {
			if useAttr && IsSubSlice(cur.Attr, attrs) {
				res = append(res, cur)
			}
			if !useAttr {
				res = append(res, cur)
			}
		}

		for c := cur.FirstChild; c != nil; c = c.NextSibling {
			q = append(q, c)
		}
	}

	return res
}

func FindTagByRegExp(doc string, regStr string) [][]string {
	r := regexp.MustCompile(regStr)
	return r.FindAllStringSubmatch(doc, -1)
}

func GetAttrVal(node *html.Node, name string) string {
	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}

	return ""
}
