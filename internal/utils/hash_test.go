package utils

import "testing"

func TestGetMD5(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Good",
			args: args{
				text: "goodtest",
			},
			want: "e5c270f25ab83e0c9b61ff6bfffc596d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMD5(tt.args.text); got != tt.want {
				t.Errorf("GetMD5() = %v, want %v", got, tt.want)
			}
		})
	}
}
