package data

type GetStrings struct {
	DefaultVisitor
	Pattern string
	Result  []GetStringsResults
}

type GetStringsResults struct {
	Value *KeyNode
}

func (visitor *GetStrings) OnKey(node *KeyNode) {
	switch node.Value.(type) {
	case string:
		visitor.Result = append(visitor.Result, GetStringsResults{Value: node})

	}
}
