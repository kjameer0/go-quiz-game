package main

import "testing"

func Test_getFileText(t *testing.T) {
	type args struct {
		userPath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFileText(tt.args.userPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFileText() = %v, want %v", got, tt.want)
			}
		})
	}
}
