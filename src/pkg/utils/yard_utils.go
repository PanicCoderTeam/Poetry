package utils

import "fmt"

func GetYardBay(block string, bayNo int) string {
	if len(block) == 0 || bayNo <= 0 {
		return ""
	}
	yardBlockBay := fmt.Sprintf("%s%03d", block, bayNo)
	return yardBlockBay
}
