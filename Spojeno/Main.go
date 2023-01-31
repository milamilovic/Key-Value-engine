package main

import (
	"Strukture/SSTable"
	"Strukture/SkipList"
)

func main() {
	sl := SkipList.MakeSkipList(10)
	sl.Add("1", []byte("a"))
	sl.Add("2", []byte("a"))
	sl.Add("3", []byte("a"))
	sl.Add("4", []byte("a"))
	sl.Add("5", []byte("a"))
	SSTable.MakeSSTable(sl.GetElements(), 1, 0)
}
