package main

type Cursor struct {
	x int
	y int
}

func NewCursor(x, y int) *Cursor {
	return &Cursor{x, y}
}

func (cursor *Cursor) Set(x, y int) {
	cursor.x = x
	cursor.y = y
}

func (cursor *Cursor) Left() {
	if cursor.x >= 0 {
		cursor.x -= 1
	}
}

func (cursor *Cursor) Right() {
	cursor.x += 1
}

func (cursor *Cursor) Up() {
	if cursor.y >= 0 {
		cursor.y -= 1
	}
}

func (cursor *Cursor) Down() {
	cursor.y += 1
}
