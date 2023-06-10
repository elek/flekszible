package data

type Set struct {
	DefaultVisitor
	Path     Path
	NewValue interface{}
}

func (visitor *Set) OnKey(node *KeyNode) {
	if visitor.Path.Match(node.Path) {
		node.Value = visitor.NewValue
	}
}

func (visitor *Set) BeforeMap(node *MapNode) {
	if visitor.Path.Parent().Equal(node.Path) {
		if node.Get(visitor.Path.Last()) == nil {
			node.PutValue(visitor.Path.Last(), visitor.NewValue)
		}
	}
}
