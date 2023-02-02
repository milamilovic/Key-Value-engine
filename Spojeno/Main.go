package main

import (
	"Strukture/BloomFilter"
	"Strukture/Cache"
	"Strukture/CountMinSketch"
	"Strukture/HyperLogLog"
	"Strukture/MemTableBTree"
	"Strukture/MemTableSkipList"
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
	cache         *Cache.Cache
	wal           *Wal.Wal
	konfiguracije map[string]int
	cms           *CountMinSketch.CountMinSketch
	hll           *HyperLogLog.HLL
	mems          *MemTableSkipList.MemTable
	memb          *MemTableBTree.MemTable
	da_li_je_skip bool
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
	engine.konfiguracije["memtable_da_li_je_skip"] = 1
}

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
func initialize(odabran string) *Engine {
	engine := Engine{}
	engine.konfiguracije = make(map[string]int)
	file, err := os.ReadFile("Data/Konfiguracije/konfiguracije.txt")
	broj := -5
	if err != nil {
		default_konfig(&engine)
	} else {
		delovi := SplitLines(string(file))
		engine.konfiguracije["memtable_max_velicina"], _ = strconv.Atoi(delovi[0])
		engine.konfiguracije["cache_size"], _ = strconv.Atoi(delovi[1])
		engine.konfiguracije["sst_level"], _ = strconv.Atoi(delovi[2])
		engine.konfiguracije["sst_index"], _ = strconv.Atoi(delovi[3])
		engine.konfiguracije["lsm_level"], _ = strconv.Atoi(delovi[4])
		engine.konfiguracije["wal_low_water_mark"], _ = strconv.Atoi(delovi[5])
		engine.konfiguracije["token_key"], _ = strconv.Atoi(delovi[6])
		engine.konfiguracije["token_maxtok"], _ = strconv.Atoi(delovi[7])
		engine.konfiguracije["token_interval"], _ = strconv.Atoi(delovi[8])
		broj, _ = strconv.Atoi(delovi[9])
	}
	if broj == 1 {
		engine.mems = MemTableSkipList.KreirajMemTable(engine.konfiguracije["memtable_max_velicina"], engine.konfiguracije["memtable_max_velicina"])
		engine.memb = nil
		engine.da_li_je_skip = true
	} else {

		engine.memb = MemTableBTree.KreirajMemTable(engine.konfiguracije["memtable_max_velicina"], engine.konfiguracije["memtable_max_velicina"])
		engine.mems = nil
		engine.da_li_je_skip = false
	}
	engine.bloom = BloomFilter.New_bloom(engine.konfiguracije["memtable_max_velicina"], float64(0.1))
	engine.cache = Cache.KreirajCache(engine.konfiguracije["cache_size"])
	engine.wal = Wal.NapraviWal("Data\\Wal", engine.konfiguracije["wal_low_water_mark"])
	return &engine
}

func main() {
	odabran := odabirMemTable()
	engine := initialize(odabran)
	menu(engine)
}

func nabavi_vrednosti_dodavanje() (string, []byte) {
	fmt.Println("Unesite vrednost kljuca: ")
	r := bufio.NewReader(os.Stdin)
	kljuc, _ := r.ReadString('\n')
	kljuc = strings.Replace(kljuc, "\n", "", 1)
	kljuc = strings.Replace(kljuc, "\r", "", 1)
	r = bufio.NewReader(os.Stdin)
	fmt.Println("Unesite vrednost pod kljucem: ")
	vrednost, _ := r.ReadString('\n')
	vrednost = strings.Replace(vrednost, "\n", "", 1)
	vrednost = strings.Replace(vrednost, "\r", "", 1)
	return kljuc, []byte(vrednost)
}

func nabavi_vrednosti_brisanje() string {
	fmt.Println("Unesite vrednost kljuca: ")
	r := bufio.NewReader(os.Stdin)
	kljuc, _ := r.ReadString('\n')
	kljuc = strings.Replace(kljuc, "\n", "", 1)
	kljuc = strings.Replace(kljuc, "\r", "", 1)
	return kljuc
}

func odabirMemTable() string {
	fmt.Println("Unesite da li za MemTable hocete da koristite Btree ili SkipList")
	fmt.Println("1:SkipList")
	fmt.Println("2:Btree")
	r := bufio.NewReader(os.Stdin)
	unos, _ := r.ReadString('\n')
	unos = strings.Replace(unos, "\n", "", 1)
	unos = strings.Replace(unos, "\r", "", 1)
	for unos != "1" && unos != "2" {
		fmt.Println("Pogresan unos!")
		r := bufio.NewReader(os.Stdin)
		unos, _ := r.ReadString('\n')
		unos = strings.Replace(unos, "\n", "", 1)
		unos = strings.Replace(unos, "\r", "", 1)
	}
	return unos

}
func menu(engine *Engine) {
	b := true
	for b == true {
		fmt.Println("****MENI*****")
		fmt.Println("1:PUT")
		fmt.Println("2:GET")
		fmt.Println("3:DELETE")
		fmt.Println("4:LIST")
		fmt.Println("5:RANGE SCAN")
		fmt.Println("6:10+")
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
		case "6":
			desetPlusMeni(engine)
			break
		case "x":
			b = false
			break
		default:
			fmt.Println("Pogresan unos")
		}

	}

}

func desetPlusMeni(engine *Engine) {
	b := true
	for b == true {
		fmt.Println("**10+ meni***")
		fmt.Println("1:NAPRAVI CMS")
		fmt.Println("2:DODAJ U CMS")
		fmt.Println("3:PROVERI CMS")
		fmt.Println("4:NAPRAVI HLL")
		fmt.Println("5:DODAJ U HLL")
		fmt.Println("6:SERIJALIZUJ HLL")
		fmt.Println("7:SIM HASH DEMONSTRACIJA")
		fmt.Println("x:povratak na obican meni")
		fmt.Println("Unesite broj ispred zeljene opcije")
		r := bufio.NewReader(os.Stdin)
		unos, _ := r.ReadString('\n')
		unos = strings.Replace(unos, "\n", "", 1)
		unos = strings.Replace(unos, "\r", "", 1)
		switch unos {
		case "1":
			makeCms(engine)
			break
		case "2":
			if engine.cms == nil {
				fmt.Println("Morate prvo napraviti cms!")
				break
			}
			key, value := nabavi_vrednosti_dodavanje()
			addCms(key, value, engine)
			break
		case "3":
			if engine.cms == nil {
				fmt.Println("Morate prvo napraviti cms!")
				break
			}
			kljuc := nabavi_vrednosti_brisanje()
			fmt.Println(checkCms(kljuc, engine))
			break
		case "4":
			makeHll(engine)
			break
		case "5":
			if engine.hll == nil {
				fmt.Println("Morate prvo napraviti hll!")
				break
			}
			key := nabavi_vrednosti_brisanje()
			addHll(key, engine)
			break
		case "6":
			if engine.hll == nil {
				fmt.Println("Morate prvo napraviti hll!")
				break
			}
			saveHll(engine)
			break
		case "7":
			makeSimHash()
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
	fmt.Println("fingerprint prvog teksta: ")
	sim1 := SimHash.SimHash("Strukture\\SimHash\\simHash.txt")
	fmt.Println(sim1)
	fmt.Println("fingerprint drugog teksta: ")
	sim2 := SimHash.SimHash("Strukture\\SimHash\\simHash2.txt")
	fmt.Println(sim2)
	fmt.Println("Hemingova razdaljina ova dva teksta: ")
	fmt.Println(SimHash.Hamming(sim1, sim2))
}

func makeCms(engine *Engine) {
	engine.cms = CountMinSketch.CreateCMS(0.1, 0.1)
}

func addCms(key string, value []byte, engine *Engine) {
	uspesno := engine.cms.Add(key, engine.cms.Hashes, int(engine.cms.M))
	if uspesno {
		fmt.Println("Element je uspesno dodat u cms!")
	}
}

func checkCms(key string, engine *Engine) int {
	return engine.cms.Cms(key, engine.cms.Hashes, int(engine.cms.M))
}

func addHll(key string, engine *Engine) {
	engine.hll.Add(key)
	fmt.Println("Element je uspesno dodat u hll!")
}

func saveHll(engine *Engine) {
	podaci := HyperLogLog.Serijalizacija(engine.hll)
	os.WriteFile("Spojeno\\Strukture\\HyperLogLog\\hll.bin", podaci, os.FileMode(os.O_RDWR))
}

func estimateHll(engine *Engine) {
	fmt.Println("Procena broja elemenata u hll: ")
	fmt.Println(engine.hll.Estimate())
}

func makeHll(engine *Engine) {
	hll := HyperLogLog.MakeHyperLogLog(HyperLogLog.HLL_MAX_PRECISION)
	engine.hll = hll
}
