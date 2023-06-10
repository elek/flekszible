package data

type GetAll struct {
	DefaultVisitor
	Path   Path
	Result []GetAllResult
}

type GetAllResult struct {
	Path  Path
	Value Node
}

func (visitor *GetAll) OnKey(node *KeyNode) {
	if visitor.Path.Match(node.Path) {
		visitor.Result = append(visitor.Result, GetAllResult{Path: node.Path, Value: node})
	}
}

func (visitor *GetAll) BeforeList(node *ListNode) {
	if visitor.Path.Match(node.Path) {
		visitor.Result = append(visitor.Result, GetAllResult{Path: node.Path, Value: node})

	}
}

func (visitor *GetAll) BeforeMap(node *MapNode) {
	if visitor.Path.Match(node.Path) {
		visitor.Result = append(visitor.Result, GetAllResult{Path: node.Path, Value: node})
	}
}
