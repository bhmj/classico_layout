package types

type RowType struct {
	Pieces  []byte
	L, M, S int
}

func NewRowType(attempt []byte) RowType {
	rt := RowType{}
	for _, v := range attempt {
		switch v {
		case 'l':
			rt.L++
			rt.Pieces = append(rt.Pieces, v)
		case 'm':
			rt.M++
			rt.Pieces = append(rt.Pieces, v)
		case 's':
			rt.S++
			rt.Pieces = append(rt.Pieces, v)
		}
	}
	return rt
}
