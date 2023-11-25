package frontier

import (
	"fmt"

	"github.com/wholesome-ghoul/web-crawler-prototype/config"
)

type PriorityQueue struct {
	last *Node
}

type Node struct {
	next  *Node
	prev  *Node
	value config.SeedUrl
}

func (p *PriorityQueue) Push(seedUrl config.SeedUrl) {
	node := Node{
		next:  nil,
		value: seedUrl,
	}

	if p.Empty() {
		node.prev = nil
		p.last = &node
		return
	}

	node.prev = p.last
	p.last.next = &node
	p.last = &node
}

func (p *PriorityQueue) Pop() *Node {
	if p.last == nil {
		return nil
	}

	node := p.last
	p.last = p.last.prev

	if p.last != nil {
		p.last.next = nil
	}

	return node
}

func (p *PriorityQueue) Empty() bool {
	return p.last == nil
}

func (p *PriorityQueue) Print() {
	curr := p.last

	for curr != nil {
		fmt.Println(curr.value)
		curr = curr.prev
	}
}

func (n *Node) Url() string {
	return n.value.Url
}

func (n *Node) Priority() uint8 {
	return n.value.Priority
}

func (n *Node) Hostname() string {
	return n.value.Hostname()
}
