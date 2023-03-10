package SkipList

import (
	"fmt"
	"math/rand"
	"time"
)

type Node struct {
	key       string
	value     []byte
	tombstone bool
}

func MakeNode(key string, value []byte, tombstone bool) *Node {
	return &Node{key, value, tombstone}

}

type SkipList struct {
	maxHeight int
	height    int
	size      int
	head      *SkipListNode
}
type SkipListNode struct {
	key       string
	value     []byte
	next      []*SkipListNode
	tombstone bool
	timestamp int64
}

func (node *SkipListNode) GetKey() string {
	return node.key
}
func (node *SkipListNode) GetValue() []byte {
	return node.value
}
func (node *SkipListNode) GetTombstone() bool {
	return node.tombstone
}
func (node *SkipListNode) GetTimeStamp() int64 {
	return node.timestamp
}
func (list *SkipList) GetElements() []*SkipListNode {
	niz := make([]*SkipListNode, 0)
	curr := list.head
	curr = curr.next[0]
	niz = append(niz, curr)
	for i := 1; i < list.size; i++ {
		curr = curr.next[0]
		niz = append(niz, curr)

	}
	return niz
}

func NapraviSkipList(maxHeight int) *SkipList {
	head := SkipListNode{key: "", value: nil, next: make([]*SkipListNode, maxHeight+1), tombstone: false, timestamp: 0}
	maxH := maxHeight
	h := 1
	size := 0
	return &SkipList{maxH, h, size, &head}
}

func (skipList *SkipList) NadjiElement(key string) (bool, *SkipListNode) {
	trenutni := skipList.head
	for i := skipList.height; i >= 0; i-- {
		sledeci := trenutni.next[i]
		for sledeci != nil {
			if sledeci.key > key {
				break
			}
			trenutni = sledeci
			sledeci = trenutni.next[i]
			if trenutni.key == key {
				if trenutni.tombstone == false {
					return true, trenutni //nasli smo ga, nije log obrisan pa vratimo true i taj elem
				} else {
					return false, trenutni //nasli smo ali je log obrisan
				}
			}
		}
	}
	return false, nil //nismo upste nasli
}

func (skipList *SkipList) Add(key string, value []byte) {
	level := skipList.Roll()
	b, pronadjeniCvor := skipList.NadjiElement(key)
	if b == false { //ili ne postoji ili je log obrisan pri cemu ga onda azuriramo
		if pronadjeniCvor != nil {
			pronadjeniCvor.tombstone = false
			pronadjeniCvor.value = value
		} else {
			now := time.Now()
			noviCvor := &SkipListNode{key, value, make([]*SkipListNode, level), false, now.Unix()}
			for i := skipList.height; i >= 0; i-- {
				trenutni := skipList.head
				sledeci := trenutni.next[i]
				for sledeci != nil {
					if sledeci == nil || sledeci.key > key {
						break
					}
					trenutni = sledeci
					sledeci = trenutni.next[i]
				}
				if i < level {
					noviCvor.next[i] = sledeci
					trenutni.next[i] = noviCvor
				}
			}
			skipList.size++
		}
	} else {
		//fmt.Println("Postoji uneti kljuc")
	}

}

func (skipList *SkipList) LogBrisanje(key string) {
	b, elem := skipList.NadjiElement(key)
	if b == true {
		elem.tombstone = true
	} else {
		if elem != nil {
			fmt.Println("Element sa unetim kljucem je vec obrisan")
		} else {
			fmt.Println("Ne postoji element sa unetim kljucem")
		}
	}

}
func (s *SkipList) Roll() int {
	level := 1
	for ; rand.Int31n(2) == 1; level++ {
		if level >= s.maxHeight {
			if level > s.height {
				s.height = level
			}
			return level
		}

		if level > s.height {
			s.height = level
		}
	}
	return level
}
