package reader_test

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/masakurapa/go-cover/internal/profile"
	"github.com/masakurapa/go-cover/internal/reader"
)

func TestReader(t *testing.T) {
	type args struct {
		pkgs profile.Packages
		path string
	}

	r := reader.New()

	baseDir, err := filepath.Abs(".")
	if err != nil {
		t.Error(err)
	}

	// read want file
	path, err := filepath.Abs("./reader.go")
	if err != nil {
		t.Error(err)
	}
	want, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name    string
		r       reader.Reader
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ファイルが読み込めること",
			r:    r,
			args: args{
				pkgs: profile.Packages{
					"github.com/masakurapa/go-cover/internal/reader": {
						Dir: baseDir,
					},
				},
				path: "github.com/masakurapa/go-cover/internal/reader/reader.go",
			},
			want:    want,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.New()
			got, err := r.Read(tt.args.pkgs, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("reader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reader.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
