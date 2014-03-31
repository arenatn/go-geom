package geom

import (
	"math"
)

type LineString struct {
	layout     Layout
	stride     int
	flatCoords []float64
}

var _ T = &LineString{}

func NewLineString(layout Layout) *LineString {
	return &LineString{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: nil,
	}
}

func NewLineStringFlat(layout Layout, flatCoords []float64) *LineString {
	return &LineString{
		layout:     layout,
		stride:     layout.Stride(),
		flatCoords: flatCoords,
	}
}

func (ls *LineString) Bounds() *Bounds {
	return NewBounds().extendFlatCoords(ls.flatCoords, 0, len(ls.flatCoords), ls.stride)
}

func (ls *LineString) Clone() *LineString {
	flatCoords := make([]float64, len(ls.flatCoords))
	copy(flatCoords, ls.flatCoords)
	return &LineString{
		layout:     ls.layout,
		stride:     ls.stride,
		flatCoords: flatCoords,
	}
}

func (ls *LineString) Coord(i int) []float64 {
	return ls.flatCoords[i*ls.stride : (i+1)*ls.stride]
}

func (ls *LineString) Coords() [][]float64 {
	return inflate1(ls.flatCoords, 0, len(ls.flatCoords), ls.stride)
}

func (ls *LineString) Ends() []int {
	return nil
}

func (ls *LineString) Endss() [][]int {
	return nil
}

func (ls *LineString) FlatCoords() []float64 {
	return ls.flatCoords
}

func (ls *LineString) LastCoord() []float64 {
	if len(ls.flatCoords) == 0 {
		return nil
	} else {
		return ls.flatCoords[len(ls.flatCoords)-ls.stride:]
	}
}

func (ls *LineString) Layout() Layout {
	return ls.layout
}

func (ls *LineString) Length() float64 {
	length := 0.0
	for i := ls.stride; i < len(ls.flatCoords); i += ls.stride {
		dx := ls.flatCoords[i] - ls.flatCoords[i-ls.stride]
		dy := ls.flatCoords[i+1] - ls.flatCoords[i+1-ls.stride]
		length += math.Sqrt(dx*dx + dy*dy)
	}
	return length
}

func (ls *LineString) Interpolate(val float64, dim int) (int, float64) {
	n := len(ls.flatCoords)
	if n == 0 {
		panic("geom: empty linestring")
	}
	if val <= ls.flatCoords[dim] {
		return 0, 0
	}
	if ls.flatCoords[n-ls.stride+dim] <= val {
		return (n - 1) / ls.stride, 0
	}
	low := 0
	high := n / ls.stride
	for low < high {
		mid := (low + high) / 2
		if val < ls.flatCoords[mid*ls.stride+dim] {
			high = mid
		} else {
			low = mid + 1
		}
	}
	low--
	val0 := ls.flatCoords[low*ls.stride+dim]
	if val == val0 {
		return low, 0
	}
	val1 := ls.flatCoords[(low+1)*ls.stride+dim]
	return low, (val - val0) / (val1 - val0)
}

func (ls *LineString) NumCoords() int {
	return len(ls.flatCoords) / ls.stride
}

func (ls *LineString) Push(coord0 []float64) error {
	if len(coord0) != ls.stride {
		return ErrStrideMismatch{Got: len(coord0), Want: ls.stride}
	}
	ls.flatCoords = append(ls.flatCoords, coord0...)
	return nil
}

func (ls *LineString) SetCoords(coords1 [][]float64) error {
	var err error
	if ls.flatCoords, err = deflate1(nil, coords1, ls.stride); err != nil {
		return err
	}
	return nil
}

func (ls *LineString) Stride() int {
	return ls.stride
}

func (ls *LineString) SubLineString(start, stop int) *LineString {
	return &LineString{
		layout:     ls.layout,
		stride:     ls.stride,
		flatCoords: ls.flatCoords[start*ls.stride : stop*ls.stride],
	}
}