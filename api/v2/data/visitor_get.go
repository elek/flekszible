package data

type Get struct {
	DefaultVisitor
	Path        Path
	ReturnValue Node
	Found       bool
}

func (visitor *Get) ValueAsString() string {
	return visitor.ReturnValue.(*KeyNode).Value.(string)
}
func (visitor *Get) OnKey(node *KeyNode) {
	if visitor.Path.Match(node.Path) && !visitor.Found {
		visitor.ReturnValue = node
		visitor.Found = true
	}
}

func (visitor *Get) BeforeList(node *ListNode) {
	if visitor.Path.Match(node.Path) && !visitor.Found {
		visitor.ReturnValue = node
		visitor.Found = true
	}
}
func (visitor *Get) BeforeMap(node *MapNode) {
	if visitor.Path.Match(node.Path) && !visitor.Found {
		visitor.ReturnValue = node
		visitor.Found = true
	}
}
