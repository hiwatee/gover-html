package reader_test

import (
	"reflect"
	"testing"

	"github.com/masakurapa/go-cover/internal/reader"
)

func TestReader(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		r       reader.Reader
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := reader.New()
			got, err := r.Read(tt.args.path)
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

func BenchmarkReader(b *testing.B) {
	r := reader.New()
	path := "github.com/masakurapa/go-cover/test/_example/example.go"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Read(path)
	}
}
