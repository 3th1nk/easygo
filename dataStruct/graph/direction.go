package graph

// Direction 连线方向
//go:generate stringer -type Direction -linecomment
type Direction int

const (
	LeftToRight Direction = 0 // -->
	RightToLeft Direction = 1 // <--
	BothWay     Direction = 2 // <->
)
