package data

type SmartGetAll struct {
	DefaultVisitor
	Path   Path
	Result []GetAllResult
}

func (visitor *SmartGetAll) OnKey(node *KeyNode) {
	if visitor.Path.Match(node.Path) {
		visitor.Result = append(visitor.Result, GetAllResult{Path: node.Path, Value: node})
	}
}

func (visitor *SmartGetAll) BeforeList(node *ListNode) {
	if visitor.Path.Match(node.Path) {
		visitor.Result = append(visitor.Result, GetAllResult{Path: node.Path, Value: node})

	}
}

var mapChildren = []Path{
	NewPath("metadata"),
	NewPath("metadata", "annotations"),
	NewPath("metadata", "labels"),
	NewPath("spec", "template", "metadata", "labels"),
	NewPath("spec", "template", "metadata"),
	NewPath("spec", "template", "metadata", "annotations"),
}

var listChildren = []Path{
	NewPath("spec", "template", "spec", "containers"),
	NewPath("spec", "template", "spec", "initContainers"),
	NewPath("spec", "template", "spec", "volumes"),
	NewPath("spec", "template", "spec", ".*ontainers", ".*", "env"),
	NewPath("spec", "template", "spec", ".*ontainers", ".*", "envFrom"),
	NewPath("spec", "template", "spec", ".*ontainers", ".*", "volumeMounts"),
	NewPath("spec", "containers"),
	NewPath("spec", "initContainers"),
	NewPath("spec", "volumes"),
	NewPath("spec", ".*ontainers", ".*", "env"),
	NewPath("spec", ".*ontainers", ".*", "envFrom"),
	NewPath("spec", ".*ontainers", ".*", "volumeMounts"),
}

func (visitor *SmartGetAll) BeforeMap(node *MapNode) {
	if visitor.Path.Match(node.Path) {
		visitor.Result = append(visitor.Result, GetAllResult{Path: node.Path, Value: node})
		return
	}
	if match, nextSegment := visitor.Path.MatchLimited(node.Path); match {
		if !node.HasKey(nextSegment) {
			for _, path := range mapChildren {
				if path.Match(node.Path.Extend(nextSegment)) {
					node.CreateMap(nextSegment)
					return
				}
			}
			for _, path := range listChildren {
				if path.Match(node.Path.Extend(nextSegment)) {
					node.CreateList(nextSegment)
					return
				}
			}
		}
	}

}
