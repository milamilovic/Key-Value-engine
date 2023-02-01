package main

import (
	"Strukture/BloomFilter"
	"Strukture/CountMinSketch"
	"Strukture/MemTable"
	"Strukture/HyperLogLog"
	"fmt"
)

type Engine struct {
	bloom    BloomFilter.BloomFilter
	memtable MemTable.memTable
}

func init() {

}

func main() {

}

func makeCms() {
	cms := CountMinSketch.CreateCMS(0.1, 0.1)
	fmt.Println(cms)
}

func makeHll() {
	hll := HyperLogLog.makeHyperLogLog(8)
	fmt.Println(hll)
}
