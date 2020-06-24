package profile_test

import (
	"testing"

	"github.com/masakurapa/go-cover/internal/profile"
)

func TestBlocks_Coverage(t *testing.T) {
	tests := []struct {
		name   string
		blocks profile.Blocks
		want   float64
	}{
		{
			name: "全レコードのcount > 1の場合100が返却される",
			blocks: profile.Blocks{
				{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 1},
				{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 2},
			},
			want: 100,
		},
		{
			name: "一部レコードがcount > 1の場合小数点を含む値が返却される",
			blocks: profile.Blocks{
				{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 1},
				{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 0},
			},
			want: 49.39759036144578, // 41 / (41 + 42) * 100
		},
		{
			name: "count > 1のレコードが無い場合0が返却される",
			blocks: profile.Blocks{
				{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 0},
				{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 0},
			},
			want: 0,
		},
		{
			name:   "Blocksが空の場合0が返却される",
			blocks: profile.Blocks{},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.blocks.Coverage(); got != tt.want {
				t.Errorf("Blocks.Coverage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkBlocks_Coverage(b *testing.B) {
	blocks := profile.Blocks{
		{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 1},
		{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blocks.Coverage()
	}
}
