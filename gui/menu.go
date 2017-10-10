package gui

// menu and menu related types
// exports nothing

import (
	mgl "github.com/go-gl/mathgl/mgl64"
)

type menu struct {
	items       []menuItem
	activeIndex int
}

type menuItem struct {
	label   string
	getIcon func() Icon
	bound   bound
}

type bound struct {
	from mgl.Vec2 // left top coord
	to   mgl.Vec2 // right bottom coord
}

func (b *bound) hits(x, y float64) bool {
	return x >= b.from.X() && x < b.to.X() &&
		y >= b.from.Y() && y < b.to.Y()
}

func newMenu() menu {
	var menu menu
	menu.items = make([]menuItem, 0)
	menu.activeIndex = -1
	return menu
}

func (m *menu) addItem(getIcon func() Icon, label string) int {
	m.items = append(m.items, menuItem{label, getIcon, bound{}})
	return len(m.items) - 1
}

func (m *menu) itemIcons() []string {
	var icons = make([]string, len(m.items))
	for i, item := range m.items {
		icons[i] = string(item.getIcon())
		i++
	}
	return icons
}

func (m *menu) setBounds(widths []int, height int) {
	from := mgl.Vec2{0, 0}
	to := mgl.Vec2{0, float64(height)}
	for i := range m.items {
		to[0] += float64(widths[i])
		m.items[i].bound = bound{from, to}
		from[0] += float64(widths[i])
	}
}
