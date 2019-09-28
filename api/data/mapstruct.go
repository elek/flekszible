package data

import (
	"fmt"
	"github.com/elek/flekszible/api/yaml"
	"strconv"
)

type Visitor interface {
	OnKey(*KeyNode)

	BeforeMap(node *MapNode)
	AfterMap(node *MapNode)
	BeforeMapItem(node *MapNode, key string, index int)
	AfterMapItem(node *MapNode, key string, index int)

	BeforeList(node *ListNode)
	AfterList(node *ListNode)
	BeforeListItem(node *ListNode, item Node, index int)
	AfterListItem(node *ListNode, item Node, index int)
}

type Node interface {
	Accept(Visitor)
}

// ----------------- KEY NODE --------------
type KeyNode struct {
	Value interface{}
	Path  Path
}

func NewKeyNode(path Path, value interface{}) KeyNode {
	return KeyNode{
		Value: value,
		Path:  path,
	}
}

func (node *KeyNode) Accept(v Visitor) {
	v.OnKey(node)
}

// ----------------- MAP NODE --------------
type MapNode struct {
	keys     []string
	children map[string]Node
	Path     Path
}

func NewMapNode(path Path) MapNode {
	m := MapNode{
		Path:     path,
		children: make(map[string]Node),
	}
	return m
}
func (node *MapNode) Put(key string, value Node) {
	node.children[key] = value
	for _, indexedKey := range node.keys {
		if indexedKey == key {
			//the key is already indexed
			return
		}
	}
	node.keys = append(node.keys, key)
}

func (node *MapNode) Accept(v Visitor) {
	v.BeforeMap(node)
	idx := 0
	for _, key := range node.keys {
		v.BeforeMapItem(node, key, idx)
		idx = idx + 1
		value := node.children[key]
		value.Accept(v)
		v.AfterMapItem(node, key, idx)

	}
	v.AfterMap(node)
}

func (node *MapNode) ToString() (string, error) {
	converted := ConvertToYaml(node)
	bytes, err := yaml.Marshal(converted)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (node MapNode) Get(s string) Node {
	result := node.children[s]
	return result
}

func (node MapNode) HasKey(s string) bool {
	_, ok := node.children[s]
	return ok
}

func (node *MapNode) Len() int {
	if node.children == nil {
		return 0
	} else {
		return len(node.children)
	}
}
func (node *MapNode) Keys() []string {
	return node.keys
}
func (node *MapNode) PutValue(key string, value interface{}) {
	node.Put(key, &KeyNode{Path: node.Path.Extend(key), Value: value})
}
func (node *MapNode) CreateMap(key string) *MapNode {
	mapNode := NewMapNode(node.Path.Extend(key))
	node.Put(key, &mapNode)
	return &mapNode
}

func (node *MapNode) CreateList(key string) *ListNode {
	listNode := NewListNode(node.Path.Extend(key))
	node.Put(key, &listNode)
	return &listNode
}

func (node *MapNode) Remove(key string) {
	delete(node.children, key)
	newKeys := make([]string, 0)
	for _, indexedKey := range node.keys {
		if indexedKey != key {
			newKeys = append(newKeys, indexedKey)
		}
	}
	node.keys = newKeys
}

func (node *MapNode) GetStringValue(s string) string {
	return node.Get(s).(*KeyNode).Value.(string)
}

// ----------------- LIST NODE --------------
type ListNode struct {
	Children []Node
	Path     Path
}

func NewListNode(path Path) ListNode {
	l := ListNode{
		Children: make([]Node, 0),
		Path:     path,
	}
	return l
}

func (node *ListNode) Append(value Node) {
	node.Children = append(node.Children, value)
}

func (node *ListNode) Accept(v Visitor) {
	v.BeforeList(node)
	for ix, value := range node.Children {
		v.BeforeListItem(node, value, ix)
		value.Accept(v)
		v.AfterListItem(node, value, ix)
	}
	v.AfterList(node)

}
func (node *ListNode) Len() int {
	return len(node.Children)
}

func (node *ListNode) AddValue(value string) {
	node.Append(&KeyNode{Path: node.Path.Extend(strconv.Itoa(len(node.Children))), Value: value})
}

func (node *ListNode) CreateMap() *MapNode {
	mapNode := NewMapNode(node.Path.Extend(strconv.Itoa(len(node.Children))))
	node.Children = append(node.Children, &mapNode)
	return &mapNode
}

// ----------------- VISITORS --------------

type PrintVisitor struct {
	DefaultVisitor
}

func (PrintVisitor) OnKey(node *KeyNode) {
	fmt.Printf("%s %s\n", node.Path.ToString(), node.Value)
}
func (PrintVisitor) BeforeMap(node *MapNode) {
	fmt.Printf("%s [map]\n", node.Path.ToString())
}
func (PrintVisitor) BeforeList(node *ListNode) {
	fmt.Printf("%s [list]\n", node.Path.ToString())
}

type DefaultVisitor struct{}

func (DefaultVisitor) OnKey(*KeyNode)                                      {}
func (DefaultVisitor) BeforeMap(node *MapNode)                             {}
func (DefaultVisitor) AfterMap(node *MapNode)                              {}
func (DefaultVisitor) BeforeMapItem(node *MapNode, key string, index int)  {}
func (DefaultVisitor) AfterMapItem(node *MapNode, key string, index int)   {}
func (DefaultVisitor) BeforeList(node *ListNode)                           {}
func (DefaultVisitor) AfterList(node *ListNode)                            {}
func (DefaultVisitor) BeforeListItem(node *ListNode, item Node, index int) {}
func (DefaultVisitor) AfterListItem(node *ListNode, item Node, index int)  {}

type Apply struct {
	DefaultVisitor
	Path     Path
	Function func(interface{}) interface{}
}

func (visitor *Apply) OnKey(node *KeyNode) {
	if visitor.Path.Match(node.Path) {
		node.Value = visitor.Function(node.Value)
	}
}

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

type GetKeys struct {
	DefaultVisitor
	Result []GetAllResult
}

type GetKeysResult struct {
	Path  Path
	Value Node
}

func (visitor *GetKeys) OnKey(node *KeyNode) {
	visitor.Result = append(visitor.Result, GetAllResult{Path: node.Path, Value: node})
}

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

type Yamlize struct {
	DefaultVisitor
	Path       Path
	Serialize  bool
	parsed     bool
	parsedPath Path
}

func (visitor *Yamlize) BeforeMapItem(node *MapNode, key string, index int) {
	if !visitor.Serialize {
		//deserialize phase
		if visitor.parsed {
			return
		}
		if match, _ := visitor.Path.MatchLimited(node.Path.Extend(key)); match {
			switch value := node.Get(key).(type) {
			case *KeyNode:
				yamlDoc := yaml.MapSlice{}

				content := value.Value.(string)
				err := yaml.Unmarshal([]byte(content), &yamlDoc)
				if err != nil {
					panic(err)
				}

				newnode, err := ConvertToNode(yamlDoc, node.Path.Extend(key))
				if err != nil {
					panic(err)
				}
				node.Put(key, newnode)
				visitor.parsed = true
				visitor.parsedPath = node.Path.Extend(key)
				break

			}
		}
	} else {
		if node.Path.Extend(key).Equal(visitor.parsedPath) {
			content, err := node.Get(key).(*MapNode).ToString()
			if err != nil {
				panic(err);
			}
			node.Put(key, &KeyNode{content, node.Path.Extend(key)})
		}
	}

}

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

type FixPath struct {
	DefaultVisitor
	CurrentPath Path
}

func (visitor *FixPath) OnKey(node *KeyNode) {
	node.Path = visitor.CurrentPath
}

func (visitor *FixPath) BeforeMap(node *MapNode) {
	node.Path = visitor.CurrentPath
}
func (visitor *FixPath) AfterMap(node *MapNode) {

}
func (visitor *FixPath) BeforeList(node *ListNode) {
	node.Path = visitor.CurrentPath
}
func (visitor *FixPath) AfterList(node *ListNode) {

}

func (visitor *FixPath) BeforeMapItem(node *MapNode, key string, index int) {
	visitor.CurrentPath = visitor.CurrentPath.Extend(key)
}
func (visitor *FixPath) AfterMapItem(node *MapNode, key string, index int) {
	visitor.CurrentPath = visitor.CurrentPath.Parent()
}
func (visitor *FixPath) BeforeListItem(node *ListNode, item Node, index int) {
	subKeyName := strconv.Itoa(index)
	if mapItem, convertable := item.(*MapNode); convertable {
		if mapItem.HasKey("name") {
			name := mapItem.Get("name").(*KeyNode).Value.(string)
			subKeyName = name
		}
	}
	visitor.CurrentPath = visitor.CurrentPath.Extend(subKeyName)
}

func (visitor *FixPath) AfterListItem(node *ListNode, item Node, index int) {
	visitor.CurrentPath = visitor.CurrentPath.Parent()
}
