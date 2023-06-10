package data

import "fmt"

type ReSet struct {
	DefaultVisitor
	Path     Path
	NewValue interface{}
}

func (visitor *ReSet) OnKey(node *KeyNode) {
	if visitor.Path.Match(node.Path) {
		if fmt.Sprintf("%s", node.Value) != "" {
			node.Value = visitor.NewValue
		}
	}
}

func (visitor *ReSet) BeforeMap(node *MapNode) {
	if visitor.Path.Parent().Equal(node.Path) {
		if node.Get(visitor.Path.Last()) == nil {
			if node.HasKey(visitor.Path.Last()) {
				node.PutValue(visitor.Path.Last(), visitor.NewValue)
			}
		}
	}
}
