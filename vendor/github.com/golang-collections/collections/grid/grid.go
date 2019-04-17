package grid

import (
	. "github.com/badgerodon/collections"
)

type (
	Grid struct {
		values     []interface{}
		cols, rows int
	}
)

func New(cols, rows int) *Grid {
	return &Grid{
		values: make([]interface{}, cols*rows),
		cols:   cols,
		rows:   rows,
	}
}

func (this *Grid) Do(f func(p Point, value interface{})) {
	for x := 0; x < this.cols; x++ {
		for y := 0; y < this.rows; y++ {
			f(Point{x, y}, this.values[x*this.cols+y])
		}
	}
}

func (this *Grid) Get(p Point) interface{} {
	if p.X < 0 || p.Y < 0 || p.X >= this.cols || p.Y >= this.rows {
		return nil
	}
	v, _ := this.values[p.X*this.cols+p.Y]
	return v
}

func (this *Grid) Rows() int {
	return this.rows
}

func (this *Grid) Cols() int {
	return this.cols
}

func (this *Grid) Len() int {
	return this.rows * this.cols
}

func (this *Grid) Set(p Point, v interface{}) {

}
