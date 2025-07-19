package utils

import (
	"fmt"
	"testing"
)

func Test_ConvertChinses(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "好雨知时节"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("%+v", ConvertChinsesSimplified2T(tt.name))
		})
	}
}
