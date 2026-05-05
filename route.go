package thttp

import "strings"

type PatternType int

const (
	ColonPattern PatternType = iota
	BracePattern
)

type SegmentType int

const (
	Static          SegmentType = iota
	ColonStyleParam             // :xxx 风格参数
	BraceStyleParam             // {xxx} 风格参数
	CatchAll                    // 通配段
)

type Segment struct {
	Type SegmentType
	Raw  string
	Name string
}

// ParsePath 解析路径，支持 :id / {id} / *filepath / {*filepath}
func ParsePath(pattern string) []Segment {
	// 去掉可能的首尾空格
	pattern = strings.TrimSpace(pattern)

	if pattern == "" || pattern == "/" {
		return []Segment{}
	}

	// 去除收尾的 /
	pattern = strings.Trim(pattern, "/")
	parts := strings.Split(pattern, "/")
	segs := make([]Segment, 0, len(parts))
	for _, p := range parts {
		// 连续 / 场景
		if p == "" {
			continue
		}

		// 通配 *
		if strings.HasPrefix(p, "*") {
			name := strings.TrimPrefix(p, "*")
			if name == "" {
				name = "_" // 匿名通配符
			}
			segs = append(segs, Segment{Type: CatchAll, Name: name})
			continue // 通配之后不应再有段，但你可以校验
		}

		// 大括号通配 {*name}
		if strings.HasPrefix(p, "{*") && strings.HasSuffix(p, "}") {
			name := p[2 : len(p)-1]
			segs = append(segs, Segment{Type: CatchAll, Name: name})
			continue
		}

		// 大括号参数 {name}
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			name := p[1 : len(p)-1]
			segs = append(segs, Segment{Type: BraceStyleParam, Name: name})
			continue
		}
		// 冒号参数 :name
		if strings.HasPrefix(p, ":") {
			name := strings.TrimPrefix(p, ":")
			segs = append(segs, Segment{Type: ColonStyleParam, Name: name})
			continue
		}

		// 静态文本
		segs = append(segs, Segment{Type: Static, Name: p})
	}
	return segs
}

// ToGinPath 转为 gin / httprouter 格式：:param 与 *catchall
func ToColonStylePattern(segs []Segment) string {
	parts := make([]string, 0, len(segs))
	for _, s := range segs {
		switch s.Type {
		case Static:
			parts = append(parts, s.Name)
		case ColonStyleParam, BraceStyleParam:
			parts = append(parts, ":"+s.Name)
		case CatchAll:
			parts = append(parts, "*"+s.Name)
		}
	}
	return "/" + strings.Join(parts, "/")
}

// ToMuxPath 转为 gorilla/mux 格式：{param} 与 {param:.*}
func ToBraceStylePattern(segs []Segment) string {
	parts := make([]string, 0, len(segs))
	for _, s := range segs {
		switch s.Type {
		case Static:
			parts = append(parts, s.Name)
		case ColonStyleParam, BraceStyleParam:
			parts = append(parts, "{"+s.Name+"}")
		case CatchAll:
			// mux 中通配通常用 {name:.*} 或直接 {name}
			parts = append(parts, "{"+s.Name+":.*}")
		}
	}
	return "/" + strings.Join(parts, "/")
}

func convertPattern(pattern string, typ PatternType) string {
	segs := ParsePath(pattern)
	switch typ {
	case BracePattern:
		return ToBraceStylePattern(segs)
	case ColonPattern:
		return ToColonStylePattern(segs)
	}

	return ToBraceStylePattern(segs)
}
