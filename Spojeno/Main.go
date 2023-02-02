package main

import (
	"Strukture/BloomFilter"
	"Strukture/Cache"
	"Strukture/CountMinSketch"
	"Strukture/HyperLogLog"
	"Strukture/MemTable"
	"Strukture/Wal"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Engine struct {
	bloom         *BloomFilter.BloomFilter
	memtable      *MemTable.MemTable
	cache         *Cache.Cache
	wal           *Wal.Wal
	konfiguracije map[string]int
}

func default_konfig(engine *Engine) {
	engine.konfiguracije["memtable_max_velicina"] = 2
	engine.konfiguracije["cache_size"] = 10
	engine.konfiguracije["sst_level"] = 3
	engine.konfiguracije["sst_index"] = 4
	engine.konfiguracije["lsm_level"] = 4
	engine.konfiguracije["wal_low_water_mark"] = 3
	engine.konfiguracije["token_key"] = 99999999
	engine.konfiguracije["token_maxtok"] = 5
	engine.konfiguracije["token_interval"] = 60
}

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
func initialize() *Engine {
	engine := Engine{}
	engine.konfiguracije = make(map[string]int)
	file, err := os.ReadFile("Data/Konfiguracije/konfiguracije.txt")
	if err != nil {
		default_konfig(&engine)
	} else {
		fmt.Println(string(file))
		delovi := SplitLines(string(file))
		fmt.Println(delovi)
		engine.konfiguracije["memtable_max_velicina"], _ = strconv.Atoi(delovi[0])
		engine.konfiguracije["cache_size"], _ = strconv.Atoi(delovi[1])
		engine.konfiguracije["sst_level"], _ = strconv.Atoi(delovi[2])
		engine.konfiguracije["sst_index"], _ = strconv.Atoi(delovi[3])
		engine.konfiguracije["lsm_level"], _ = strconv.Atoi(delovi[4])
		engine.konfiguracije["wal_low_water_mark"], _ = strconv.Atoi(delovi[5])
		engine.konfiguracije["token_key"], _ = strconv.Atoi(delovi[6])
		engine.konfiguracije["token_maxtok"], _ = strconv.Atoi(delovi[7])
		engine.konfiguracije["token_interval"], _ = strconv.Atoi(delovi[8])
	}
<<<<<<< HEAD
	engine.bloom = BloomFilter.New_bloom(engine.konfiguracije["memtable_max_velicina"], float64(0.1))
	engine.memtable = MemTable.KreirajMemTable(engine.konfiguracije["memtable_max_velicina"], engine.konfiguracije["memtable_max_velicina"])
	engine.cache = Cache.KreirajCache(engine.konfiguracije["cache_size"])
	engine.wal = Wal.NapraviWal("Data\\Wal", engine.konfiguracije["wal_low_water_mark"])
=======
	engine.bloom = BloomFilter.New_bloom(engine.konfiguracije["memtable_max_velicina"], 0.1)
	engine.memtable = MemTable.KreirajMemTable(engine.konfiguracije["memtable_max_velicina"], engine.konfiguracije["memtable_max_velicina"])
	engine.wal = Wal.NapraviWal("Spojeno\\Data\\Wal", engine.konfiguracije["wal_low_water_mark"])
>>>>>>> 001c9ebf4e37aebc4738131d0bc6d79ae381b4f6
	return &engine
}

func main() {
	engine := initialize()
	fmt.Println()
	fmt.Println()
	fmt.Println(engine)
	makeCms()
	makeHll()
}

func makeCms() {
	cms := CountMinSketch.CreateCMS(0.1, 0.1)
	fmt.Println(cms)
}

func makeHll() {
	hll := HyperLogLog.MakeHyperLogLog(8)
	fmt.Println(hll)
}
