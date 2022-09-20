package application

import "testing"

func Test_parseInput(t *testing.T) {

	tests := []struct {
		name    string
		arg     string
		want    int
		want1   int
		want2   string
		wantErr bool
	}{
		{
			name:    "task sample",
			arg:     "/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg",
			want:    300,
			want1:   200,
			want2:   "http://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg",
			wantErr: false,
		},
		{
			name:    "short success",
			arg:     "/0/0/link",
			want:    0,
			want1:   0,
			want2:   "http://link",
			wantErr: false,
		},
		{
			name:    "not enough params",
			arg:     "/200/100",
			wantErr: true,
		},
		{
			name:    "wrong int 1",
			arg:     "/one/100/link",
			wantErr: true,
		},
		{
			name:    "wrong int 2",
			arg:     "/200/two/link",
			wantErr: true,
		},
		{
			name:    "wrong link",
			arg:     "/100/200/!@#$%^",
			wantErr: true,
		},
		{
			name:    "http only",
			arg:     "/100/200/ftp://aaa",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := parseInput(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseInput() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseInput() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("parseInput() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
