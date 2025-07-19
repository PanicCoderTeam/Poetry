package utils

import (
	"github.com/siongui/gojianfan"
)

func ConvertChinsesSimplified2T(simplified string) string {
	traditional := gojianfan.S2T(simplified) // 自动处理多音字和地区差异
	return traditional
}

func ConvertChinsesTraditional2S(traditional string) string {
	simplified := gojianfan.T2S(traditional) // 自动处理多音字和地区差异
	return simplified
}
