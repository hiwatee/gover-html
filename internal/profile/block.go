package profile

import "sort"

type Blocks []Block

type Block struct {
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
	NumStmt   int
	Count     int
}

func (blocks *Blocks) Filter() Blocks {
	newBlocks := make(Blocks, 0, len(*blocks))
	for _, b := range *blocks {
		i := newBlocks.Index(&b)

		if i == -1 {
			newBlocks = append(newBlocks, b)
			continue
		}

		if b.Count > 0 {
			newBlocks[i] = b
		}
	}

	return newBlocks
}

func (blocks *Blocks) Sort() {
	b := *blocks
	sort.SliceStable(*blocks, func(i, j int) bool {
		bi, bj := b[i], b[j]
		return bi.StartLine < bj.StartLine || bi.StartLine == bj.StartLine && bi.StartCol < bj.StartCol
	})
}

func (blocks *Blocks) Index(block *Block) int {
	for i, b := range *blocks {
		if b.StartLine == block.StartLine &&
			b.StartCol == block.StartCol &&
			b.EndLine == block.EndLine &&
			b.EndCol == block.EndCol {
			return i
		}
	}
	return -1
}

func (blocks *Blocks) Coverage() float64 {
	var total, covered int64
	for _, b := range *blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}

	if total == 0 {
		return 0
	}

	return float64(covered) / float64(total) * 100
}
