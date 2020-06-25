package html

import (
	"bufio"
	"html/template"
	"io"

	"github.com/masakurapa/go-cover/internal/profile"
)

type TemplateFile struct {
	Name     string
	Body     template.HTML
	Coverage float64
}

func genSource(w io.Writer, src []byte, prof *profile.Profile) error {
	dst := bufio.NewWriter(w)

	bi := 0
	si := 0
	line := 1
	col := 1

	for si < len(src) {
		if len(prof.Blocks) > bi {
			block := prof.Blocks[bi]
			if block.StartLine == line && block.StartCol == col {
				if block.Count == 0 {
					dst.WriteString(`<span class="uncov">`)
				} else {
					dst.WriteString(`<span class="cov">`)
				}
			}
			if block.EndLine == line && block.EndCol == col || line > block.EndLine {
				dst.WriteString(`</span>`)
				bi++
				continue
			}
		}

		b := src[si]
		switch b {
		case '<':
			dst.WriteString("&lt;")
		case '>':
			dst.WriteString("&gt;")
		case '&':
			dst.WriteString("&amp;")
		case '\t':
			dst.WriteString("        ")
		default:
			dst.WriteByte(b)
		}

		if b == '\n' {
			line++
			col = 0
		}

		si++
		col++
	}

	return dst.Flush()
}
