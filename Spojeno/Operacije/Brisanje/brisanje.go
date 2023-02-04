package brisanje

import (
	"Strukture/Cache"
	"Strukture/MemTableBTree"
	"Strukture/MemTableSkipList"
)

func Obrisi_skiplist(key string, mt *MemTableSkipList.MemTable, c *Cache.Cache) bool {
	// logicko brisanje
	obrisan := mt.BrisiElement(key)
	if obrisan {
		c.ObrisiIzCache(key)
		//kada je memTable pun ako ima neki logicki obrisan element?
		return true
	}
	return false
}

func Obrisi_bstablo(key string, mt *MemTableBTree.MemTable, c *Cache.Cache) bool {
	// logicko brisanje
	obrisan := mt.BrisiElement(key)
	if obrisan {
		c.ObrisiIzCache(key)
		return true
	}
	return false
}
