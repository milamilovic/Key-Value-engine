package brisanje

import (
	"Strukture/Cache"
	"Strukture/MemTable"
)

func obrisi(key string, mt MemTable.MemTable, c Cache.Cache) bool {
	// logicko brisanje
	obrisan := mt.BrisiElement(key)
	if obrisan {
		c.ObrisiIzCache(key)
		return true
	}
	return false
}
