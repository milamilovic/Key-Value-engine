// package main

// import (
// 	"math/rand"
// )

// func main() {
// 	skipList := makeSkipList(15)
// 	skipList.add("3", []byte("a"))
// 	skipList.add("5", []byte("a"))
// 	skipList.add("9", []byte("a"))
// 	skipList.add("1", []byte("a"))
// 	skipList.add("7", []byte("a"))
// 	skipList.add("1", []byte("ponovo isti kljuc"))
// 	skipList.add("4", []byte("a"))
// 	b, _ := skipList.findElement("1")
// 	print("\n", b)
// 	b, _ = skipList.findElement("5")
// 	print("\n", b)
// 	b, _ = skipList.findElement("7")
// 	print("\n", b)
// 	b, _ = skipList.findElement("12") //ovo je false
// 	print("\n", b)

// 	skipList.logDelete("4")
// 	b, _ = skipList.findElement("4") //false jer je logicki obrisan
// 	print("\n", b, "\n")
// 	skipList.logDelete("4") //brisemo ga opet

// 	skipList.add("4", []byte("a"))   //postoji ali je obrisan pa ga azuriramo
// 	b, _ = skipList.findElement("4") //azurirano je pa je true
// 	print("\n", b)

// }

// type SkipList struct {
// 	maxHeight int
// 	height    int
// 	size      int
// 	head      *SkipListNode
// }
// type SkipListNode struct {
// 	key       string
// 	value     []byte
// 	next      []*SkipListNode
// 	tombstone bool
// }

// func makeSkipList(maxHeight int) *SkipList {
// 	head := SkipListNode{key: "", value: nil, next: make([]*SkipListNode, maxHeight+1), tombstone: false}
// 	maxH := maxHeight
// 	h := 1
// 	size := 0
// 	return &SkipList{maxH, h, size, &head}
// }

// func (skipList *SkipList) findElement(key string) (bool, *SkipListNode) {
// 	trenutni := skipList.head
// 	for i := skipList.height; i >= 0; i-- {
// 		sledeci := trenutni.next[i]
// 		for sledeci != nil {
// 			if sledeci.key > key {
// 				break
// 			}
// 			trenutni = sledeci
// 			sledeci = trenutni.next[i]
// 			if trenutni.key == key {
// 				if trenutni.tombstone == false {
// 					return true, trenutni //nasli smo ga, nije log obrisan pa vratimo true i taj elem
// 				} else {
// 					return false, trenutni //nasli smo ali je log obrisan
// 				}
// 			}
// 		}
// 	}
// 	return false, nil //nismo upste nasli
// }

// func (skipList *SkipList) add(key string, value []byte) {
// 	level := skipList.roll()
// 	b, pronadjeniCvor := skipList.findElement(key)
// 	if b == false { //ili ne postoji ili je log obrisan pri cemu ga onda azuriramo
// 		if pronadjeniCvor != nil {
// 			pronadjeniCvor.tombstone = false
// 			pronadjeniCvor.value = value
// 		} else {
// 			noviCvor := &SkipListNode{key, value, make([]*SkipListNode, level), false}
// 			for i := skipList.height; i >= 0; i-- {
// 				trenutni := skipList.head
// 				sledeci := trenutni.next[i]
// 				for sledeci != nil {
// 					if sledeci == nil || sledeci.key > key {
// 						break
// 					}
// 					trenutni = sledeci
// 					sledeci = trenutni.next[i]
// 				}
// 				if i < level {
// 					noviCvor.next[i] = sledeci
// 					trenutni.next[i] = noviCvor
// 					skipList.size++
// 				}
// 			}
// 		}
// 	} else {
// 		print("Postoji vec dati kljuc")
// 	}

// }

// func (skipList *SkipList) logDelete(key string) {
// 	b, elem := skipList.findElement(key)
// 	if b == true {
// 		elem.tombstone = true
// 	} else {
// 		if elem != nil {
// 			print("Element sa unetim kljucem je vec obrisan")
// 		} else {
// 			print("Ne postoji element sa unetim kljucem")
// 		}
// 	}

// }
// func (s *SkipList) roll() int {
// 	level := 1
// 	for ; rand.Int31n(2) == 1; level++ {
// 		if level >= s.maxHeight {
// 			if level > s.height {
// 				s.height = level
// 			}
// 			return level
// 		}

// 		if level > s.height {
// 			s.height = level
// 		}
// 	}
// 	return level
// }
