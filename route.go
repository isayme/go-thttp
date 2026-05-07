package thttp

import "strings"

type SegmentType int

const (
	Static   SegmentType = iota // static
	Param                       // param
	CatchAll                    // wildcard
)

type Segment struct {
	Type SegmentType
	Raw  string
	Name string
}

// ParsePath parse pattern
func ParsePath(pattern string) []Segment {
	pattern = strings.TrimSpace(pattern)

	if pattern == "" || pattern == "/" {
		return []Segment{}
	}

	pattern = strings.Trim(pattern, "/")
	parts := strings.Split(pattern, "/")
	segs := make([]Segment, 0, len(parts))
	for _, p := range parts {
		// case: double slash
		if p == "" {
			continue
		}

		// catch-all *
		if strings.HasPrefix(p, "*") {
			name := strings.TrimPrefix(p, "*")
			if name == "" {
				panic("catch-all should be named")
			}
			segs = append(segs, Segment{Type: CatchAll, Name: name, Raw: p})
			continue
		}

		// std catch-all
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "...}") {
			name := p[1 : len(p)-4]
			if name == "" {
				panic("catch-all should be named")
			}
			segs = append(segs, Segment{Type: CatchAll, Name: name, Raw: p})
			continue
		}

		// catch-all {*name}
		if strings.HasPrefix(p, "{*") && strings.HasSuffix(p, "}") {
			name := p[2 : len(p)-1]
			segs = append(segs, Segment{Type: CatchAll, Name: name, Raw: p})
			continue
		}

		// param style {name}
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			name := p[1 : len(p)-1]
			segs = append(segs, Segment{Type: Param, Name: name, Raw: p})
			continue
		}

		// pram style :name
		if strings.HasPrefix(p, ":") {
			name := strings.TrimPrefix(p, ":")
			segs = append(segs, Segment{Type: Param, Name: name, Raw: p})
			continue
		}

		// static
		segs = append(segs, Segment{Type: Static, Name: p, Raw: p})
	}
	return segs
}
