package MemTable

import (
	"Strukture/SSTable"
	"Strukture/SkipList"
	"fmt"
)

type memTable struct {
	elementi         *SkipList.SkipList
	velicina         int //velicina skipListe
	maxVelicina      int //maksimalna velicina za memTable
	trenutnaVelicina int
}

func CreateMemTable(max, velicina int) *memTable {
	elementi := SkipList.MakeSkipList(velicina)
	return &memTable{elementi, velicina, max, 0}
}

func (memTable *memTable) Add(key string, value []byte) {
	b, cvor := memTable.elementi.FindElement(key)
	if b == false {
		if cvor == nil {
			memTable.elementi.Add(key, value)
			fmt.Println("Ubacili smo novi element u skip listu")
			memTable.trenutnaVelicina++
		}
	}
}
func (memTable *memTable) Update(key string, value []byte) {
	b, cvor := memTable.elementi.FindElement(key)
	if b == true { //nasao je elemnt i menja mu value
		memTable.elementi.Add(key, value)
		fmt.Println("Izmenili smo element u skip listi")
	} else {
		if cvor != nil { //cvor je logicki obrisan
			memTable.elementi.Add(key, value) //izmenice mu i tombstone na false
		}
	}
}
func (memTable *memTable) DeleteElement(key string) {
	b, cvor := memTable.elementi.FindElement(key)
	if b == true { //nasao je elemnt i menja mu value
		memTable.elementi.LogDelete(key)
		fmt.Println("Izbrisali smo element u skip listi")
	} else {
		if cvor != nil { //cvor je logicki obrisan
			fmt.Println("Element je vec logicki obrisan")
		} else {
			fmt.Println("Nema element sa unetim kljucem")
		}
	}
}

func (memTable *memTable) CheckFlush() bool {
	if memTable.maxVelicina <= memTable.trenutnaVelicina {
		return true //treba flush odraditi
	} else {
		return false
	}
}

var i int = 0

func (memTable *memTable) Flush() {
	memTable.WriteSSTable(i)
	memTable = CreateMemTable(15, 20) //pre ovoga treba upisati na disk, SStable
}

func (memTable *memTable) WriteSSTable(i int) {
	i++
	SSTable.MakeSSTable(memTable.elementi.GetElements(), 1, i)
}
