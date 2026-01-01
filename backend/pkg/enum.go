package pkg

import (
	"fmt"
)

// 生成性别枚举并导出
type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

func (g Gender) String() string { return string(g) }
func (g Gender) ZhString() string {
	switch g {
	case Male:
		return "男"
	case Female:
		return "女"
	default:
		return "未知"
	}
}

func (g Gender) IsValid() bool {
	switch g {
	case Male, Female:
		return true
	default:
		return false
	}
}

func ParseGender(s string) (Gender, error) {
	g := Gender(s)
	if !g.IsValid() {
		return "", fmt.Errorf("invalid Gender: %s", s)
	}
	return g, nil
}
