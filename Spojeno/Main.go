package main

import (
	"Strukture/BloomFilter"
	"Strukture/Cache"
	"Strukture/CountMinSketch"
	"Strukture/MemTable"
	"Strukture/HyperLogLog"
	"Strukture/Wal"
	"fmt"
	"os"
)

type Engine struct {
	bloom         BloomFilter.BloomFilter
	memtable      MemTable.MemTable
	cache         Cache.Cache
	wal           Wal.Wal
	konfiguracije map[string]int
}

func default_konfig(engine *Engine) {
	engine.konfiguracije["memtable_max_velicina"] = 15
	engine.konfiguracije["cache_size"] = 10
	engine.konfiguracije["sst_level"] = 3
	engine.konfiguracije["sst_index"] = 4
	engine.konfiguracije["lsm_level"] = 4
	engine.konfiguracije["wal_low_water_mark"] = 3
	engine.konfiguracije["token_key"] = 99999999
	engine.konfiguracije["token_maxtok"] = 5
	engine.konfiguracije["token_interval"] = 60
}

func initialize() *Engine {
	engine := Engine{}
	engine.konfiguracije = make(map[string]int)
	file, err := os.ReadFile("Data/Konfiguracije/konfiguracije.txt")
	if err != nil {
		default_konfig(&engine)
	} else {
		fmt.Println(string(file))
	}
	engine.bloom = BloomFilter.New_bloom(engine.konfiguracije["memtable_max_velicina"], 0.1)
	engine.memtable = MemTable.CreateMemTable(engine.konfiguracije["memtable_max_velicina"], engine.konfiguracije["memtable_max_velicina"])
	engine.wal = Wal.NapraviWal("Spojeno\\Data\\Wal", engine.konfiguracije["wal_low_water_mark"])
	return &engine
}

func main() {
	engine := initialize()
	fmt.Println(engine)
}

func makeCms() {
	cms := CountMinSketch.CreateCMS(0.1, 0.1)
	fmt.Println(cms)
}

func makeHll() {
	hll := HyperLogLog.MakeHyperLogLog(8)
	fmt.Println(hll)
}
