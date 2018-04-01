package gui

// menu and menu related types
// exports nothing

import (
	mgl "github.com/go-gl/mathgl/mgl64"
)

type menuItem struct {
	label      string
	getIcon    func() Icon
	isDisabled func() bool
	bound      bound
}
type menu struct {
	items       []menuItem
	activeIndex int
	showTooltip bool
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

func (m *menu) addItem(getIcon func() Icon, isDisabled func() bool, label string) int {
	m.items = append(m.items, menuItem{label, getIcon, isDisabled, bound{}})
	return len(m.items) - 1
}

func (m *menu) itemIcons() []string {
	var icons = make([]string, len(m.items))
	for i, item := range m.items {
		icons[i] = string(item.getIcon())
	}
	return icons
}

func (m *menu) tooltip() string {
	if m.activeIndex > -1 && m.showTooltip {
		return m.items[m.activeIndex].label
	} else {
		return ""
	}
}

func (m *menu) disabled() []bool {
	var disabled = make([]bool, len(m.items))
	for i, item := range m.items {
		disabled[i] = item.isDisabled()
	}
	return disabled
}

func (m *menu) setBounds(widths []float64, height float64, top float64) {
	from := mgl.Vec2{0, top}
	to := mgl.Vec2{0, top + height}
	for i := range m.items {
		to[0] += widths[i]
		m.items[i].bound = bound{from, to}
		from[0] += widths[i]
	}
}
