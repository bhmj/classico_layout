package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/bhmj/classico_layout/internal/html"
	"github.com/bhmj/classico_layout/internal/pkg/config"
	"github.com/bhmj/classico_layout/internal/pkg/types"
)

type service struct {
	cfg    *config.Config
	matrix []types.RowType
}

type Service interface {
	GenerateMatrix()
	GetMatrix() []types.RowType
	Run(context.Context)
}

func NewService(cfg *config.Config) Service {
	return &service{cfg: cfg}
}

func (s *service) GenerateMatrix() {
	total := generateRows(nil, s.cfg.Road.Width)

	large, medium, small := s.cfg.Pallet.Large, s.cfg.Pallet.Medium, s.cfg.Pallet.Small

	// limit the number of fewest pieces per row (calculation speed optimization)
	var fewest *int = &small
	if large < *fewest {
		fewest = &large
	}
	if medium < *fewest {
		fewest = &medium
	}
	limit := 8 * s.cfg.Road.Width / (((large + medium + small) - *fewest) / *fewest) / 5
	*fewest = limit

	var matrix []types.RowType
	for n, row := range total {
		if row.L > large || row.M > medium || row.S > small {
			continue
		}
		dup := false
		for i := 0; i < n; i++ {
			if total[i].L == row.L && total[i].M == row.M && total[i].S == row.S {
				dup = true
				break
			}
		}
		if dup {
			continue
		}

		matrix = append(matrix, row)
	}

	s.matrix = matrix
}

func (s *service) GetMatrix() []types.RowType {
	return s.matrix
}

func generateRows(prefix []byte, maxLen int) []types.RowType {
	result := make([]types.RowType, 0)
	vars := []string{"l--", "m-", "s"}
	for _, v := range vars {
		attempt := append(prefix, []byte(v)...)
		switch {
		case len(attempt) == maxLen:
			result = append(result, types.NewRowType(attempt))
		case len(attempt) < maxLen:
			result = append(result, generateRows(attempt, maxLen)...)
		}
	}
	return result
}

func (s *service) Run(ctx context.Context) {
	f, err := html.CreateOutputFile("layout.html", "Layout")
	if err != nil {
		return
	}

	large, medium, small := s.cfg.Pallet.Large, s.cfg.Pallet.Medium, s.cfg.Pallet.Small
	colorClass := []string{"odd", "even"}
	var colorRemain [6]int
	for ilayer := 0; ilayer < s.cfg.Pallet.Layers; ilayer++ {
		fmt.Printf("layer %d of %d\n", ilayer, s.cfg.Pallet.Layers)
		remLarge, remMedium, remSmall, layout := generateLayout(nil, large, medium, small, s.cfg.Road.Width, s.matrix, 0)
		layout, colorRemain = paintLayout(layout, large, medium, small, remLarge, remMedium, remSmall, colorRemain)
		large, medium, small = s.cfg.Pallet.Large+remLarge, s.cfg.Pallet.Medium+remMedium, s.cfg.Pallet.Small+remSmall
		html.WriteLayout(f, shuffleLayout(layout), colorClass[ilayer%len(colorClass)], colorRemain)
		if ilayer == s.cfg.Pallet.Layers-1 {
			fmt.Printf("remainder:\n  large\t%d\n  medium\t%d\n  small\t%d\n", remLarge, remMedium, remSmall)
		}
	}
	fmt.Fprintln(f, "</body></html>")
	f.Close()

	fmt.Println("done")
}

func generateLayout(input []types.RowType, l, m, s, cols int, matrix []types.RowType, level int) (int, int, int, []types.RowType) {
	if l*3+m*2+s < cols {
		return l, m, s, nil
	}
	minRemainder := cols
	var best []types.RowType
	rl, rm, rs := 100, 100, 100
	for irow, row := range matrix {
		if level < 2 {
			fmt.Printf("%*d (%d/%d)\n", level+1, level, irow+1, len(matrix))
		}
		if row.L > l || row.M > m || row.S > s {
			continue
		}
		rlf, rmf, rsf, layout := generateLayout(input, l-row.L, m-row.M, s-row.S, cols, matrix, level+1)
		remainder := rlf*3 + rmf*2 + rsf
		if remainder < minRemainder {
			rl, rm, rs = rlf, rmf, rsf
			minRemainder = remainder
			best = append(input, row)
			best = append(best, layout...)
		}
	}
	if minRemainder < cols {
		return rl, rm, rs, best
	}
	return rl, rm, rs, nil
}

// prevRemain: [l,m,s,L,M,S]
func paintLayout(layout []types.RowType, nlarge, nmedium, nsmall int, rlarge, rmedium, rsmall int, remain [6]int) ([]types.RowType, [6]int) {
	lms := [3]int{nlarge, nmedium, nsmall}
	rlms := [3]int{rlarge, rmedium, rsmall}
	chars := [3]byte{'l', 'm', 's'}
	upChars := [3]byte{'L', 'M', 'S'}
	for size := 0; size < len(lms); size++ {
		n := lms[size] - remain[size] - remain[size+3] - rlms[size] // number of elems to paint
		index := make([]int, 0)
		for irow, row := range layout {
			for i := range row.Pieces {
				if row.Pieces[i] == chars[size] {
					index = append(index, irow*1024+i)
				}
			}
		}
		rand.Shuffle(len(index), func(i, j int) { index[i], index[j] = index[j], index[i] })
		for i := 0; i < n-n/2+remain[size+3]; i++ {
			layout[index[i]/1024].Pieces = []byte(string(layout[index[i]/1024].Pieces))
			layout[index[i]/1024].Pieces[index[i]%1024] = upChars[size]
		}

		remain[size] = rlms[size] - rlms[size]/2
		remain[size+3] = rlms[size] / 2
	}

	return layout, remain
}

func shuffleLayout(layout []types.RowType) []types.RowType {
	for r, row := range layout {
		rand.Shuffle(len(row.Pieces), func(i, j int) { row.Pieces[i], row.Pieces[j] = row.Pieces[j], row.Pieces[i] })
		layout[r].Pieces = []byte(string(row.Pieces))
	}
	rand.Shuffle(len(layout), func(i, j int) { layout[i], layout[j] = layout[j], layout[i] })
	return layout
}
