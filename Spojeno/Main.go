package main

import (
	// "Strukture/BloomFilter"
	// "Strukture/CountMinSketch"
	// "Strukture/HyperLogLog"
	"Strukture/MemTable"
	// "Strukture/SkipList"
)

func main() {
	memTable := MemTable.CreateMemTable(20, 15)
	memTable.Add("1", []byte("a"))
	memTable.Update("1", []byte("b"))
	memTable.DeleteElement("1")
	flush := memTable.CheckFlush()
	if flush == true {
		memTable.Flush()
	}
}
