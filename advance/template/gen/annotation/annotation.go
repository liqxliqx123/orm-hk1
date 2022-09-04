package annotation

import (
	"go/ast"
	"strings"
)

type Annotations[NN ast.Node] struct {
	Node NN
	Ans  []Annotation
}

func (a Annotations[NN]) Get(key string) (Annotation, bool) {
	for _, an := range a.Ans {
		if an.Key == key {
			return an, true
		}
	}
	return Annotation{}, false
}

type Annotation struct {
	Key   string
	Value string
}

func newAnnotations[NN ast.Node](n NN, cg *ast.CommentGroup) Annotations[NN] {
	if cg == nil || len(cg.List) == 0 {
		return Annotations[NN]{Node: n}
	}
	ans := make([]Annotation, 0, len(cg.List))
	for _, c := range cg.List {
		text, ok := extractContent(c)
		if !ok {
			continue
		}
		if strings.HasPrefix(text, "@") {
			segs := strings.SplitN(text, " ", 2)
			key := segs[0][1:]
			val := ""
			if len(segs) == 2 {
				val = segs[1]
			}
			ans = append(ans, Annotation{
				Key:   key,
				Value: val,
			})
		}
	}
	return Annotations[NN]{
		Node: n,
		Ans:  ans,
	}
}

func extractContent(c *ast.Comment) (string, bool) {
	text := c.Text
	if strings.HasPrefix(text, "// ") {
		return text[3:], true
	} else if strings.HasPrefix(text, "/* ") {
		length := len(text)
		return text[3 : length-2], true
	}
	return "", false
}
