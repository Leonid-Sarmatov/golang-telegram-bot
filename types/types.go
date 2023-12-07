package types

import (
	"strconv"
	"fmt"
)

// Интерфейс, описывающий обстрактный узел дерева меню
type Node interface {
	IsEnding() bool
	GetName() string
	GetCallback() string
	GetSubmenus() []Node
	GetPrice() string
}

// Структура, описывающая узел дерева меню, не являющийся конечным
type NodeMenu struct {
	Name string
	Callback string
	SubMenu []Node
}

// Структура, описывающая конечные узлы дерева, то есть позиции для заказа
type EndingMenu struct {
	Name string
	Price string
	Callback string
}



func (m NodeMenu) IsEnding() bool {
	return false
}

func (m NodeMenu) GetName() string {
	return m.Name
}

func (m NodeMenu) GetCallback() string {
	return m.Callback
}

func (m NodeMenu) GetSubmenus() []Node {
	return m.SubMenu
}

func (m NodeMenu) GetPrice() string {
	return "None"
}



func (e EndingMenu) IsEnding() bool {
	return true
}

func (e EndingMenu) GetName() string {
	return e.Name
}

func (e EndingMenu) GetCallback() string {
	return e.Callback
}

func (e EndingMenu) GetSubmenus() []Node {
	return nil
}

func (e EndingMenu) GetPrice() string {
	return e.Price
}



// Структура описывающая корзину
type Cart struct {
	SumPrice string
	PositionItems map[string][]EndingMenu
}

// Добавления позиции в корзину
func (c *Cart) AddPositionItem(e EndingMenu) error {
	d, err := strconv.Atoi(e.Price)
	if err != nil {
		return fmt.Errorf("ERROR: Price value is not integer. Order: %s\n", e.Name)
	}

	nowPrice, _ := strconv.Atoi(c.SumPrice)
	newPrice := nowPrice+d
	c.SumPrice = strconv.Itoa(newPrice)

	// Если такая позиция уже есть в заказе, значит добавляем еще один такой же товар
	if val, ok := c.PositionItems[e.Name]; ok && len(val) != 0 {
		c.PositionItems[e.Name] = append(c.PositionItems[e.Name], e)
	// Иначе просто добавляем позицию
	} else {
		c.PositionItems[e.Name] = []EndingMenu{e}
	}
	return nil
}

// Удаление позиции из корзины
func (c *Cart) RemovePositionItem(e EndingMenu) error {
	d, err := strconv.Atoi(e.Price)
	if err != nil {
		return fmt.Errorf("ERROR: Price value is not integer. Order: %s\n", e.Name)
	}

	// Если такие товары уже лежат в корзине, то уменьшаем их количество
	if val, ok := c.PositionItems[e.Name]; ok {
		if len(val) != 0 {
			c.PositionItems[e.Name] = c.PositionItems[e.Name][1:]
		} else {
			delete(c.PositionItems, e.Name)
		}
		// Уменьшаем счетчик цены
		nowPrice, _ := strconv.Atoi(c.SumPrice)
		newPrice := nowPrice-d
		c.SumPrice = strconv.Itoa(newPrice)
	// Иначе просто убираем позицию
	} else {
		return fmt.Errorf("ERROR: Can not remove position, cart is empty. Order: %s\n", e.Name)
	}
	return nil
}
