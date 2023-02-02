package Cache

import (
	"fmt"
)

type DoubleLinkedList struct {
	head    *CacheNode
	tail    *CacheNode
	dllSize int
}

type CacheNode struct {
	key   string
	value []byte
	prev  *CacheNode
	next  *CacheNode
}

type Cache struct {
	list        *DoubleLinkedList
	hashes      map[string][]byte
	numberOfNew int
	maxSize     int
}

func NapraviListu(size int) *DoubleLinkedList {
	return &DoubleLinkedList{head: nil, tail: nil, dllSize: size}
}
func NapraviCvor(key string, value []byte) *CacheNode {
	node := CacheNode{key, value, nil, nil}
	return &node
}
func (dll *DoubleLinkedList) DodajCvor(key string, value []byte) bool {
	newNode := NapraviCvor(key, value)
	if dll.head != nil {
		newNode.next = dll.head
		dll.head.prev = newNode
		dll.head = newNode

	} else {
		dll.head = newNode
		return true
	}
	if dll.tail == nil {
		dll.tail = dll.head.next
	}
	return true
}

func (dll *DoubleLinkedList) NadjiUListi(key string) *CacheNode {
	current := dll.head
	for current != nil {
		if current.key == key {
			return current
		}
		current = current.next
	}
	return nil
}
func (dll *DoubleLinkedList) ObrisiPoslednjiCvor() string {
	key := dll.tail.key
	dll.tail = dll.tail.prev
	dll.tail.next = nil
	return key
}

func (dll *DoubleLinkedList) ObrisiCvor(node *CacheNode) {
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		dll.tail = node.prev
	}
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		dll.head = node.next
	}

}
func (dll *DoubleLinkedList) PrintList() {
	current := dll.head
	for current != nil {
		fmt.Print(current.key, " ")
		current = current.next
	}
}

func KreirajCache(size int) *Cache {
	dll := NapraviListu(size)
	return &Cache{dll, make(map[string][]byte), 0, size}
}

func (cache *Cache) DodajUCache(key string, value []byte) {
	_, ok := cache.NadjiUCache(key)
	if ok {
		node := cache.list.NadjiUListi(key) //ako postoji vec u cache trazimo ga cvor u listi i brisemo
		cache.list.ObrisiCvor(node)
		delete(cache.hashes, key) //brisemo staru vrednost iz hashes
		b := cache.list.DodajCvor(key, value)
		if b {
			fmt.Println("Uspesno dodavanje")
			cache.hashes[key] = value
		}
	} else {
		if len(cache.hashes) >= cache.maxSize {
			del := cache.list.ObrisiPoslednjiCvor() //brise se poslednji i njegova vrednost u hashes tabeli
			delete(cache.hashes, del)
		}
		b := cache.list.DodajCvor(key, value)
		if b {
			fmt.Println("Uspesno dodavanje elementa")
			cache.hashes[key] = value
			cache.numberOfNew++
		}
	}
}
func (cache *Cache) NadjiUCache(key string) ([]byte, bool) {
	val, ok := cache.hashes[key]
	if ok {
		return val, true
	} else {
		return nil, false
	}
}
func (cache *Cache) GetFromCache(key string) *CacheNode {
	_, ok := cache.hashes[key]
	if ok {
		node := cache.list.NadjiUListi(key)
		cache.list.ObrisiCvor(node)
		b := cache.list.DodajCvor(node.key, node.value) //pristupili smo mu pa ga brisemo, i ubacujemo na pocetak
		if b {
			fmt.Println("Uspesan pristup elementu")
			return node
		} else {
			return nil
		}

	} else {
		return nil
	}
}
func (cache *Cache) ObrisiIzCache(key string) bool {
	_, ok := cache.hashes[key]
	if ok {
		node := cache.list.NadjiUListi(key)
		cache.list.ObrisiCvor(node) //prvo izbrisemo cvor
		delete(cache.hashes, key)   // zatim izbrisemo iz hashes
		return true
	} else {
		//nema unetog kljuca u hases, znaci nema tog cvora
		return false
	}
}
func (cache *Cache) Print() {
	cache.list.PrintList()
}
