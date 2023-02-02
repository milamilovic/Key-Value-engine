package brisanje

import (
	"Strukture/Cache"
	"Strukture/MemTableBTree"
	"Strukture/MemTableSkipList"
)

func obrisi_skiplist(key string, mt MemTableSkipList.MemTable, c Cache.Cache) bool {
	// logicko brisanje
	obrisan := mt.BrisiElement(key)
	if obrisan {
		c.ObrisiIzCache(key)
		return true
	}
	return false
}

func obrisi_bstablo(key string, mt MemTableBTree.MemTable, c Cache.Cache) bool {
	// logicko brisanje
	obrisan := mt.BrisiElement(key)
	if obrisan {
		c.ObrisiIzCache(key)
		return true
	}
	return false
}
