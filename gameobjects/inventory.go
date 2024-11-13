package gameobjects

import rl "github.com/gen2brain/raylib-go/raylib"

type ItemType int

const (
	Weapon ItemType = iota
	HealthPack
	Other
)

type Item struct {
	Type  ItemType
	Name  string
	Image rl.Texture2D
}

type Inventory struct {
	Slots        []Item
	MaxSlots     int
	IsOpen       bool
	SelectedSlot int
}

func NewInventory(maxSlots int) Inventory {
	return Inventory{
		Slots:    make([]Item, maxSlots),
		MaxSlots: maxSlots,
	}
}

// Method to add an item to the inventory
func (inv *Inventory) AddItem(item Item) bool {
	for i := 0; i < inv.MaxSlots; i++ {
		if inv.Slots[i].Type == Other { 
			inv.Slots[i] = item
			return true
		}
	}
	return false
}
