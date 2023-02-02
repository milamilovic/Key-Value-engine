package MemTableSkipList

import (
	"Strukture/SSTable"
	"Strukture/SkipList"
	"fmt"
)

type MemTable struct {
	elementi         *SkipList.SkipList
	velicina         int //velicina skipListe
	maxVelicina      int //maksimalna velicina za memTable
	trenutnaVelicina int
}

func (memTable *MemTable) NadjiElement(kljuc string) (bool, []byte) {
	b, cvor := memTable.elementi.NadjiElement(kljuc)
	if cvor != nil {
		value := cvor.GetValue()
		return b, value
	}
	return false, nil
}
func KreirajMemTable(max, velicina int) *MemTable {
	elementi := SkipList.NapraviSkipList(velicina)
	return &MemTable{elementi, velicina, max, 0}
}

func (memTable *MemTable) Add(key string, value []byte) {
	b, cvor := memTable.elementi.NadjiElement(key)
	if b == false {
		if cvor == nil {
			memTable.elementi.Add(key, value)
			//fmt.Println("Ubacili smo novi element u skip listu")
			memTable.trenutnaVelicina++
		}
	}
}
func (memTable *MemTable) Update(key string, value []byte) {
	b, cvor := memTable.elementi.NadjiElement(key)
	if b == true { //nasao je elemnt i menja mu value
		memTable.elementi.Add(key, value)
		fmt.Println("Izmenili smo element u skip listi")
	} else {
		if cvor != nil { //cvor je logicki obrisan
			memTable.elementi.Add(key, value) //izmenice mu i tombstone na false
		}
	}
}
func (memTable *MemTable) BrisiElement(key string) bool {
	b, cvor := memTable.elementi.NadjiElement(key)
	if b == true { //nasao je elemnt i menja mu value
		memTable.elementi.LogBrisanje(key)
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
}

func (memTable *MemTable) NapraviSSTable(i int) {
	i++
	SSTable.NapraviSSTable(memTable.elementi.GetElements(), 1, i)
}
