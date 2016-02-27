package mango

import (
	"fmt"
	"strconv"
	"strings"
)

type mediaType struct {
	mainType    string
	subType     string
	rangeParams string
	q           float32
}

func (m mediaType) Empty() bool {
	return len(m.mainType) == 0 || len(m.subType) == 0
}

func (m mediaType) String() string {
	return m.mainType + "/" + m.subType + m.rangeParams
}

func newMediaType(s string) (*mediaType, error) {
	mt := mediaType{q: 1}
	if s == "" {
		mt.mainType = "*"
		mt.subType = "*"
		return &mt, nil
	}

	parts := strings.Split(s, ";")
	rnge := strings.Split(parts[0], "/")

	if len(rnge) != 2 {
		return nil, fmt.Errorf("invalid media type: %q", parts[0])
	}
	mt.mainType = strings.TrimSpace(rnge[0])
	mt.subType = strings.TrimSpace(rnge[1])
	// media range parameters and quality factors are
	// rarely used in Accept headers, but...
	for i := 1; i < len(parts); i++ {
		p := strings.TrimSpace(parts[i])
		s := strings.Replace(p, " ", "", 1)
		if !strings.HasPrefix(s, "q=") {
			mt.rangeParams += ";" + p
			continue
		}
		// extract quality factor and stop.
		// ignore any extension parameters
		qs := strings.Split(s, "=")
		q, _ := strconv.ParseFloat(qs[1], 32)
		mt.q = float32(q)
		break
	}

	return &mt, nil
}

type mediaTypes []mediaType

func (m mediaTypes) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m mediaTypes) Len() int      { return len(m) }
func (m mediaTypes) Less(i, j int) bool {
	if m[i].q != m[j].q {
		return m[i].q > m[j].q
	}
	if m[i].mainType == m[j].mainType {
		if m[i].subType == m[j].subType {
			return len(m[i].rangeParams) > len(m[j].rangeParams)
		}
		return m[i].subType > m[j].subType
	}
	return m[i].mainType > m[j].mainType
}
