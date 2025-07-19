package utils

import "testing"

func TestGetYardBay(t *testing.T) {
	type args struct {
		block string
		bayNo int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test",
			args: args{
				block: "A01",
				bayNo: 1,
			},
			want: "A01001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetYardBay(tt.args.block, tt.args.bayNo); got != tt.want {
				t.Errorf("GetYardBay() = %v, want %v", got, tt.want)
			}
		})
	}
}
