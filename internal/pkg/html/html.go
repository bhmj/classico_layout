package html

import (
	"fmt"
	"os"

	"github.com/bhmj/classico_layout/internal/pkg/types"
)

func WriteHTML(matrix []types.RowType, fname string, title string) {
	f, _ := CreateOutputFile(fname, title)
	defer f.Close()

	WriteLayout(f, matrix, "", [6]int{})
	fmt.Fprintln(f, "</body></html>")
}

func CreateOutputFile(fname string, title string) (*os.File, error) {
	f, err := os.Create(fname)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(f, "<html><head><title>%s</title><style>\n.odd { background-color: #EEFFEE; }\n.even { background-color: #FFEEEE; }</style></head><body>", title)
	return f, nil
}

func WriteLayout(f *os.File, layout []types.RowType, colorClass string, colorRemain [6]int) {
	fmt.Fprintf(f, `<div class="%s">`, colorClass)
	for _, row := range layout {
		for _, elem := range row.Pieces {
			var name string
			switch elem {
			case 'l':
				name = "large"
			case 'm':
				name = "medium"
			case 's':
				name = "small"
			case 'L':
				name = "large2"
			case 'M':
				name = "medium2"
			case 'S':
				name = "small2"
			}
			fmt.Fprintf(f, `<img src="img/%s.png">`, name)
		}
		fmt.Fprintln(f, "<br>")
	}
	fmt.Fprintf(f, "</div>\n")
	if colorClass != "" {
		fmt.Fprintf(f, `<div style="margin: 10px 0 20px;">`)
		imgs := []string{"large", "medium", "small", "large2", "medium2", "small"}
		for i := 0; i < 6; i++ {
			if colorRemain[i] > 0 {
				fmt.Fprintf(f, `<img src="img/%s.png"> = %d&nbsp;&nbsp;&nbsp;`, imgs[i], colorRemain[i])
			}
		}
		fmt.Fprintf(f, "</div>")
	}
}
