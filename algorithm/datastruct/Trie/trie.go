package main

import "fmt"

func NewRuneTrie() *RuneTrie {
	return &RuneTrie{root: &runeTrieNode{child: make(map[rune]*runeTrieNode)}}
}

type RuneTrie struct {
	root *runeTrieNode
}

type runeTrieNode struct {
	child map[rune]*runeTrieNode
	Value interface{}
}

func (r *RuneTrie) Insert(key string, value interface{}) {
	node := r.root
	for _, c := range key {
		if n, ok := node.child[c]; !ok {
			newNode := &runeTrieNode{child: make(map[rune]*runeTrieNode)}
			node.child[c] = newNode
			node = newNode
		} else {
			node = n
		}
	}
	node.Value = value
}

func (r *RuneTrie) Get(key string) interface{} {
	node := r.root
	for _, c := range key {
		if n, ok := node.child[c]; !ok {
			return nil
		} else {
			node = n
		}
	}
	return node.Value
}

func main() {
	t := NewRuneTrie()
	t.Insert("hello", "hello")
	t.Insert("world", "world")
	t.Insert("河北", "河北")
	t.Insert("湖南", "湖南")
	t.Insert("湖北", "湖北省")
	fmt.Println(t.Get("湖北"))
}
