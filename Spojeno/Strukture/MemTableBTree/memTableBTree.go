package MemTableBTree

import (
	"Strukture/BTree"
	"Strukture/SSTable"
	"Strukture/SkipList"
	"fmt"
)

var Niz []*BTree.Node
var nizSize int

type MemTable struct {
	Elementi         *BTree.Btree
	velicina         int //velicina skipListe
	maxVelicina      int //maksimalna velicina za memTable
	trenutnaVelicina int
}

func KreirajMemTable(max, velicina int) *MemTable {
	elementi := BTree.MakeBtree(velicina)
	return &MemTable{elementi, velicina, max, 0}
}
func (memTable *MemTable) NadjiElement(kljuc string) (bool, []byte) {
	b, cvor := memTable.Elementi.FindElement(kljuc)
	if cvor != nil {
		value := cvor.GetValue()
		return b, value
	}
	return false, nil
}

func (memTable *MemTable) Add(key string, value []byte) {
	b, cvor := memTable.Elementi.FindElement(key)
	if b == false {
		if cvor == nil {
			Niz = append(Niz, BTree.MakeNode(key, value, false))
			nizSize++
			memTable.Elementi.Add(key, value)
			fmt.Println("Ubacili smo novi element u skip listu")
			memTable.trenutnaVelicina++
		}
	}
}
func (memTable *MemTable) Update(key string, value []byte) {
	b, cvor := memTable.Elementi.FindElement(key)
	if b == true { //nasao je elemnt i menja mu value
		memTable.Elementi.Add(key, value)
		fmt.Println("Izmenili smo element u skip listi")
	} else {
		if cvor != nil { //cvor je logicki obrisan
			memTable.Elementi.Add(key, value) //izmenice mu i tombstone na false
		}
	}
}
func (memTable *MemTable) BrisiElement(key string) bool {
	b, cvor := memTable.Elementi.FindElement(key)
	if b == true { //nasao je elemnt i menja mu value
		memTable.Elementi.LogDel(key)
		fmt.Println("Izbrisali smo element u skip listi")
		return true
	} else {
		if cvor != nil { //cvor je logicki obrisan
			fmt.Println("Element je vec logicki obrisan")
			return false
		} else {
			fmt.Println("Nema element sa unetim kljucem")
			return false
		}
	}
}

func (memTable *MemTable) ProveriFlush() bool {
	if memTable.maxVelicina <= memTable.trenutnaVelicina {
		return true //treba flush odraditi
	} else {
		return false
	}
}

var i int = 0

func (memTable *MemTable) Flush() {
	memTable.NapraviSSTable(i)
	memTable = KreirajMemTable(10, 10) //pre ovoga treba upisati na disk, SStable
}

func (memTable *MemTable) NapraviSSTable(i int) {
	sl := SkipList.NapraviSkipList(nizSize)
	for _, elem := range Niz {
		sl.Add(elem.GetKey(), elem.GetValue())
	}
	i++
	SSTable.NapraviSSTable(sl.GetElements(), 1, i)
}
