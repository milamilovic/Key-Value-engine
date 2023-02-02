package main

import (
	"Strukture/BloomFilter"
	"Strukture/Cache"
	"Strukture/CountMinSketch"
	"Strukture/HyperLogLog"
	"Strukture/MemTable"
	"Strukture/SimHash"
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
	engine.bloom = BloomFilter.New_bloom(engine.konfiguracije["memtable_max_velicina"], float64(0.1))
	engine.memtable = MemTable.KreirajMemTable(engine.konfiguracije["memtable_max_velicina"], engine.konfiguracije["memtable_max_velicina"])
	engine.cache = Cache.KreirajCache(engine.konfiguracije["cache_size"])
	engine.wal = Wal.NapraviWal("Data\\Wal", engine.konfiguracije["wal_low_water_mark"])
	return &engine
}

func main() {
	engine := initialize()
	fmt.Println()
	fmt.Println()
	fmt.Println(engine)
	menu()
	makeCms()
	makeHll()
	makeSimHash()
}
func menu() {
	b := true
	for b == true {
		fmt.Println("****MENI*****")
		fmt.Println("1:PUT")
		fmt.Println("2:GET")
		fmt.Println("3:DELETE")
		fmt.Println("4:LIST")
		fmt.Println("5:RANGE SCAN")
		fmt.Println("x:kraj programa")
		fmt.Println("Unesite broj ispred zeljene opcije")
		r := bufio.NewReader(os.Stdin)
		unos, _ := r.ReadString('\n')
		unos = strings.Replace(unos, "\n", "", 1)
		unos = strings.Replace(unos, "\r", "", 1)
		switch unos {
		case "1":
			//put()
			break
		case "2":
			//get()
			break
		case "3":
			//delete()
			break
		case "4":
			//list()
			break
		case "5":
			//rangeScan()
			break
		case "x":
			b = false
			break
		default:
			fmt.Println("Pogresan unos")
		}

	}

}

func makeSimHash() {
	sim1 := SimHash.SimHash("Strukture\\SimHash\\simHash.txt")
	fmt.Println(sim1)
	sim2 := SimHash.SimHash("Strukture\\SimHash\\simHash2.txt")
	fmt.Println(sim2)
	fmt.Println(SimHash.Hamming(sim1, sim2))
}

func makeCms() {
	cms := CountMinSketch.CreateCMS(0.1, 0.1)
	fmt.Println(cms)
}

func makeHll() {
	hll := HyperLogLog.MakeHyperLogLog(8)
	fmt.Println(hll)
}
