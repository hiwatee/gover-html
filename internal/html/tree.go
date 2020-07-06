package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"sync"

	"github.com/masakurapa/go-cover/internal/profile"
	"github.com/masakurapa/go-cover/internal/reader"
)

var (
	writePool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

type TemplateData struct {
	Tree  template.HTML
	Files []File
}

type File struct {
	ID       int
	Name     string
	Body     template.HTML
	Coverage float64
}

func WriteTreeView(
	reader reader.Reader,
	out io.Writer,
	profiles profile.Profiles,
	tree profile.Tree,
) error {
	data := TemplateData{
		Tree:  template.HTML(directoryTree(tree)),
		Files: make([]File, 0, len(profiles)),
	}

	pkgs, err := profiles.ToPackages()
	if err != nil {
		return err
	}

	for _, p := range profiles {
		b, err := reader.Read(pkgs, &p)
		if err != nil {
			return fmt.Errorf("can't read %q: %v", p.FileName, err)
		}

		buf := writePool.Get().(*bytes.Buffer)
		writeSource(buf, b, &p)

		data.Files = append(data.Files, File{
			ID:       p.ID,
			Name:     p.FileName,
			Body:     template.HTML(buf.String()),
			Coverage: p.Blocks.Coverage(),
		})

		buf.Reset()
		writePool.Put(buf)
	}

	return parsedTreeTemplate.Execute(out, data)
}

func directoryTree(tree profile.Tree) string {
	buf := writePool.Get().(*bytes.Buffer)
	writeDirectoryTree(buf, tree, 0)
	s := buf.String()
	buf.Reset()
	writePool.Put(buf)
	return s
}

func writeDirectoryTree(buf *bytes.Buffer, tree profile.Tree, indent int) {
	for _, t := range tree {
		buf.WriteString(fmt.Sprintf(
			`<div style="padding-inline-start: %dpx">%s</div>`,
			indent*30,
			t.Name,
		))

		writeDirectoryTree(buf, t.SubDirs, indent+1)

		for _, p := range t.Profiles {
			_, f := filepath.Split(p.FileName)
			buf.WriteString(fmt.Sprintf(
				`<div class="file" style="padding-inline-start: %dpx" id="tree%d" onclick="change(%d);">%s (%.1f%%)</div>`,
				(indent+1)*30,
				p.ID,
				p.ID,
				f,
				p.Blocks.Coverage(),
			))
		}
	}
}

func writeSource(buf *bytes.Buffer, src []byte, prof *profile.Profile) {
	bi := 0
	si := 0
	line := 1
	col := 1
	cov0 := false
	cov1 := false

	buf.WriteString("<ol>")

	for si < len(src) {
		if col == 1 {
			buf.WriteString("<li>")
			if cov0 {
				buf.WriteString(`<span class="cov0">`)
			}
			if cov1 {
				buf.WriteString(`<span class="cov1">`)
			}
		}

		if len(prof.Blocks) > bi {
			block := prof.Blocks[bi]
			if block.StartLine == line && block.StartCol == col {
				if block.Count == 0 {
					buf.WriteString(`<span class="cov0">`)
					cov0 = true
				} else {
					buf.WriteString(`<span class="cov1">`)
					cov1 = true
				}
			}
			if block.EndLine == line && block.EndCol == col || line > block.EndLine {
				buf.WriteString(`</span>`)
				bi++
				cov0 = false
				cov1 = false
				continue
			}
		}

		b := src[si]
		switch b {
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '&':
			buf.WriteString("&amp;")
		case '\t':
			buf.WriteString("        ")
		default:
			buf.WriteByte(b)
		}

		if b == '\n' {
			if cov0 || cov1 {
				buf.WriteString("</span>")
			}
			buf.WriteString("</li>")
			line++
			col = 0
		}

		si++
		col++
	}

	buf.WriteString("</ol>")
}
