package rocinante

import (
	"fmt"
	"strings"
)

type trie struct {
	roots map[string]*node
}

func newTrie() *trie {
	return &trie{roots: make(map[string]*node)}
}

func (t *trie) parsePattern(pattern string) []string {
	if pattern == "/" {
		return []string{"/"}
	}

	if !strings.HasPrefix(pattern, "/") {
		panic("invalid pattern")
	}

	parts := strings.Split(pattern, "/")[1:]
	lp := len(parts)
	for i, part := range parts {
		if i < lp-1 {
			if part == "" {
				panic("invalid pattern")
			}
			if strings.HasPrefix(part, "*") {
				panic(`"*" must be in the last position`)
			}
		} else {
			if part == "" {
				parts = parts[:lp-1]
			}
		}

		if strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*") {
			if len([]rune(part)) <= 1 {
				panic("invalid param")
			}
		}
	}
	return parts
}

func (t *trie) insert(pattern string, method string) {
	existsPattern, _, exists := t.search(pattern, method)
	if exists {
		panic(fmt.Sprintf(`there is a conflict between pattern:"%s" and pattern:"%s"`, pattern, existsPattern))
	}

	parts := t.parsePattern(pattern)

	root := t.roots[method]
	for i, part := range parts {
		isWild := false
		isEnd := false
		if strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*") {
			isWild = true
		}
		if isWild {
			root.wildChildPart = part
		}

		if i == len(parts)-1 {
			isEnd = true
		}

		_, ok := root.children[part]
		if !ok {
			if strings.HasPrefix(part, ":") {
				for key := range root.children {
					if (strings.HasPrefix(key, ":") || strings.HasPrefix(key, "*")) && key != part {
						panic(fmt.Sprintf(`there is a conflict between part:"%s" and part:"%s"`, part, key))
					}
				}
			}

			root.children[part] = &node{
				part:     part,
				children: make(map[string]*node),
				isWild:   isWild,
				isEnd:    isEnd,
			}
		}

		root = root.children[part]
	}
}

func (t *trie) search(pattern string, method string) (string, Params, bool) {
	_, ok := t.roots[method]
	if !ok {
		t.roots[method] = &node{children: make(map[string]*node)}
	}

	parts := t.parsePattern(pattern)
	lp := len(parts)

	params := make(Params)

	root := t.roots[method]

	sb := &strings.Builder{}

	for i, part := range parts {
		node, matched := root.children[part]
		if matched {
			if part == "/" {
				return "/", params, true
			}

			sb.WriteString(fmt.Sprintf("/%s", part))

			if node.isWild {
				params[part[1:]] = part
			}
		} else {
			//for childPart, childNode := range root.children {
			//	if strings.HasPrefix(childPart, ":") {
			//		sb.WriteString(fmt.Sprintf("/%s", childPart))
			//		params[childPart[1:]] = part
			//		matched = true
			//
			//		node = childNode
			//		break
			//	} else if strings.HasPrefix(childPart, "*") {
			//		sb.WriteString(fmt.Sprintf("/%s", childPart))
			//		wb := strings.Builder{}
			//		for _, s := range parts[i:] {
			//			wb.WriteString(fmt.Sprintf("/%s", s))
			//		}
			//		params[childPart[1:]] = wb.String()
			//		return sb.String(), params, true
			//	}
			//}

			wildChildPart := root.wildChildPart
			if wildChildPart != "" {
				node = root.children[wildChildPart]
				if strings.HasPrefix(wildChildPart, ":") {
					sb.WriteString(fmt.Sprintf("/%s", wildChildPart))
					params[wildChildPart[1:]] = part
					matched = true
				} else if strings.HasPrefix(wildChildPart, "*") {
					sb.WriteString(fmt.Sprintf("/%s", wildChildPart))
					wb := strings.Builder{}
					for _, s := range parts[i:] {
						wb.WriteString(fmt.Sprintf("/%s", s))
					}
					params[wildChildPart[1:]] = wb.String()
					return sb.String(), params, true
				}
			}

			if !matched {
				return "", nil, false
			}
		}

		if i == lp-1 {
			if node.isEnd {
				return sb.String(), params, true
			}
		}

		root = node
	}
	return "", nil, false
}

type node struct {
	part          string
	children      map[string]*node
	wildChildPart string
	isWild        bool
	isEnd         bool
}
