package main

import (
	brisanje "Operacije/Brisanje"
	citanje "Operacije/Citanje"
	"Operacije/List"
	dodavanje "Operacije/Pisanje"
	"Operacije/RangeScan"
	"Strukture/BloomFilter"
	"Strukture/Cache"
	"Strukture/CountMinSketch"
	"Strukture/HyperLogLog"
	"Strukture/MemTableBTree"
	"Strukture/MemTableSkipList"
	"Strukture/SSTable"
	"Strukture/SimHash"
	"Strukture/TokenBucket"
	"Strukture/Wal"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Engine struct {
	bloom           *BloomFilter.BloomFilter
	cache           *Cache.Cache
	wal             *Wal.Wal
	konfiguracije   map[string]int
	cms             []string
	cms_podaci      map[string][]byte
	cms_pokaz       map[string]*CountMinSketch.CountMinSketch
	hll             []string
	hll_podaci      map[string][]byte
	hll_pokaz       map[string]*HyperLogLog.HLL
	mems            *MemTableSkipList.MemTable
	memb            *MemTableBTree.MemTable
	da_li_je_skip   bool
	token           *TokenBucket.TokenBucket
	indexBloom      int
	levelBloom      int
	simHashevi      map[string]string
	simHashKljucevi []string
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
	engine.konfiguracije["da_li_je_vise_fajlova"] = 1
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
	podaci, _ := os.ReadFile("Data/Konfiguracije/podaci.txt")
	broj := -5
	if err != nil {
		default_konfig(&engine)
	} else {
		pod := SplitLines(string(podaci))
		engine.levelBloom, _ = strconv.Atoi(pod[0])
		engine.indexBloom, _ = strconv.Atoi(pod[1])
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
		engine.konfiguracije["da_li_je_vise_fajlova"], _ = strconv.Atoi(delovi[10])
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
	engine.simHashevi = make(map[string]string)
	engine.cms_podaci = make(map[string][]byte)
	engine.cms_pokaz = make(map[string]*CountMinSketch.CountMinSketch)
	engine.hll_podaci = make(map[string][]byte)
	engine.hll_pokaz = make(map[string]*HyperLogLog.HLL)
	if engine.konfiguracije["da_li_je_vise_fajlova"] == 1 {
		engine.bloom = BloomFilter.New_bloom(engine.konfiguracije["memtable_max_velicina"], float64(0.1), engine.levelBloom, engine.indexBloom)
	}
	//engine.bloom = BloomFilter.New_bloom(engine.konfiguracije["memtable_max_velicina"], float64(0.1), engine.levelBloom, engine.indexBloom)
	engine.cache = Cache.KreirajCache(engine.konfiguracije["cache_size"])
	engine.wal = Wal.NapraviWal("Data\\Wal", engine.konfiguracije["wal_low_water_mark"])
	engine.token = TokenBucket.NewTokenBucket(engine.konfiguracije["token_key"], engine.konfiguracije["token_maxtok"], int64(engine.konfiguracije["token_interval"]))
	return &engine
}

func main() {
	engine := Engine{}
	engine = *initialize()
	menu(&engine)
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

// func odabirMemTable() string {
// 	fmt.Println("Unesite da li za MemTable hocete da koristite Btree ili SkipList")
// 	fmt.Println("1:SkipList")
// 	fmt.Println("2:Btree")
// 	r := bufio.NewReader(os.Stdin)
// 	unos, _ := r.ReadString('\n')
// 	unos = strings.Replace(unos, "\n", "", 1)
// 	unos = strings.Replace(unos, "\r", "", 1)
// 	for unos != "1" && unos != "2" {
// 		fmt.Println("Pogresan unos!")
// 		r := bufio.NewReader(os.Stdin)
// 		unos, _ := r.ReadString('\n')
// 		unos = strings.Replace(unos, "\n", "", 1)
// 		unos = strings.Replace(unos, "\r", "", 1)
// 	}
// 	return unos

// }
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
			key, value := nabavi_vrednosti_dodavanje()
			if engine.da_li_je_skip {
				if engine.konfiguracije["da_li_je_vise_fajlova"] == 1 {
					dodavanje.Dodaj_skiplist(key, value, engine.mems, engine.wal, engine.token, engine.bloom, true)
				} else {
					dodavanje.Dodaj_skiplist(key, value, engine.mems, engine.wal, engine.token, engine.bloom, false)
				}
				b, _ := engine.mems.ProveriFlush()
				if b {
					engine.mems.Flush(engine.indexBloom, engine.konfiguracije["da_li_je_vise_fajlova"])
					engine.indexBloom++
					if engine.konfiguracije["da_li_je_vise_fajlova"] == 1 {
						engine.bloom = nil
						engine.bloom = BloomFilter.New_bloom(engine.konfiguracije["memtable_max_velicina"], float64(0.1), 1, engine.indexBloom)
					}
					engine.mems = MemTableSkipList.KreirajMemTable(engine.konfiguracije["memtable_max_velicina"], engine.konfiguracije["memtable_max_velicina"])

				}
			} else {
				if engine.konfiguracije["da_li_je_vise_fajlova"] == 1 {
					dodavanje.Dodaj_bstablo(key, value, engine.memb, engine.wal, engine.token, engine.bloom, true)
				} else {
					dodavanje.Dodaj_bstablo(key, value, engine.memb, engine.wal, engine.token, engine.bloom, false)
				}
				b, _ := engine.memb.ProveriFlush()
				if b {
					engine.memb.Flush(engine.indexBloom, engine.konfiguracije["da_li_je_vise_fajlova"])
					engine.indexBloom++
					if engine.konfiguracije["da_li_je_vise_fajlova"] == 1 {
						engine.bloom = nil
						engine.bloom = BloomFilter.New_bloom(engine.konfiguracije["memtable_max_velicina"], float64(0.1), 1, engine.indexBloom)
					}
					engine.memb = MemTableBTree.KreirajMemTable(engine.konfiguracije["memtable_max_velicina"], engine.konfiguracije["memtable_max_velicina"])

				}
			}
			break
		case "2":
			key := nabavi_vrednosti_brisanje()
			if engine.da_li_je_skip {
				if engine.konfiguracije["da_li_je_vise_fajlova"] == 1 {
					b, value := citanje.CitajSkip(key, engine.mems, engine.cache)
					if b {
						fmt.Println("Nasao je kljuc, vrednost je:", value)
					} else {
						fmt.Println("Nije nasao uneti kljuc")
					}
				} else {
					b, value := citanje.CitajSkipJedanFajl(key, engine.mems, engine.cache)
					if b {
						fmt.Println("Nasao je kljuc, vrednost je:", value)
					} else {
						fmt.Println("Nije nasao uneti kljuc")
					}
				}
			} else {
				if engine.konfiguracije["da_li_je_vise_fajlova"] == 1 {
					b, value := citanje.CitajBTree(key, engine.memb, engine.cache)
					if b {
						fmt.Println("Nasao je kljuc, vrednost je:", value)
					} else {
						fmt.Println("Nije nasao uneti kljuc")
					}
				} else {
					b, value := citanje.CitajBTreeJedanFajl(key, engine.memb, engine.cache)
					if b {
						fmt.Println("Nasao je kljuc, vrednost je:", value)
					} else {
						fmt.Println("Nije nasao uneti kljuc")
					}
				}
			}
			break
		case "3":
			key := nabavi_vrednosti_brisanje()
			if engine.da_li_je_skip {
				brisanje.Obrisi_skiplist(key, engine.mems, engine.cache)
			} else {
				brisanje.Obrisi_bstablo(key, engine.memb, engine.cache)
			}
			break
		case "4":
			path1, _ := filepath.Abs("../Spojeno/Data/SSTableData")
			path := strings.ReplaceAll(path1, `\`, "/")
			_, _, index_files, _, _ := citanje.Svi_fajlovi(path)
			var kljucevi []string
			for i := range index_files {
				indFile, err := os.OpenFile(path+"/"+index_files[i], os.O_RDONLY, 0666)
				if err != nil {
					panic(err)
				}
				kljucevi = append(kljucevi, SSTable.Svi_kljucevi_jednog_fajla(indFile)...)
			}
			if engine.da_li_je_skip {
				for i := range engine.mems.Elementi.GetElements() {
					if engine.mems.Elementi.GetElements()[i] != nil {
						kljucevi = append(kljucevi, engine.mems.Elementi.GetElements()[i].GetKey())
					}
				}
			} else {
				for _, elem := range MemTableBTree.Niz {
					kljucevi = append(kljucevi, elem.GetKey())
				}
			}
			fmt.Println("Unesite podstring kljuca koji trazite:")
			podstring := nabavi_vrednosti_brisanje()
			fmt.Println("Unesite velicinu stranice: ")
			r := bufio.NewReader(os.Stdin)
			unos, _ := r.ReadString('\n')
			unos = strings.Replace(unos, "\n", "", 1)
			unos = strings.Replace(unos, "\r", "", 1)
			velicina, err := strconv.ParseInt(unos, 10, 0)
			if err != nil {
				fmt.Println("Niste uneli broj!")
				break
			}
			fmt.Println("Unesite redni broj stranice koju zelite: ")
			r = bufio.NewReader(os.Stdin)
			unos, _ = r.ReadString('\n')
			unos = strings.Replace(unos, "\n", "", 1)
			unos = strings.Replace(unos, "\r", "", 1)
			redni_broj, err := strconv.ParseInt(unos, 10, 0)
			if err != nil {
				fmt.Println("Niste uneli broj!")
				break
			}
			trazeni_kljucevi := List.List(podstring, int(velicina), int(redni_broj), kljucevi)
			vrednosti := make([][]byte, len(trazeni_kljucevi))
			for i := 0; i < len(trazeni_kljucevi); i++ {
				if engine.da_li_je_skip {
					_, value := citanje.CitajSkip(trazeni_kljucevi[i], engine.mems, engine.cache)
					vrednosti[i] = value
				} else {
					_, value := citanje.CitajBTree(trazeni_kljucevi[i], engine.memb, engine.cache)
					vrednosti[i] = value
				}
			}
			fmt.Println("Vrednosti dobijene list-om su: ")
			fmt.Println(vrednosti)
			break
		case "5":
			//rangeScan()
			var kljucevi_memtable []string
			if engine.da_li_je_skip {
				for i := range engine.mems.Elementi.GetElements() {
					if engine.mems.Elementi.GetElements()[i] != nil {
						kljucevi_memtable = append(kljucevi_memtable, engine.mems.Elementi.GetElements()[i].GetKey())
					}
				}
			} else {
				for _, elem := range MemTableBTree.Niz {
					kljucevi_memtable = append(kljucevi_memtable, elem.GetKey())
				}
			}
			fmt.Println("Minimalan kljuc:")
			kljuc1 := nabavi_vrednosti_brisanje()
			fmt.Println("Maksimalan kljuc:")
			kljuc2 := nabavi_vrednosti_brisanje()
			if kljuc1 > kljuc2 {
				fmt.Println("Odnos kljuceva nije dobar!")
				break
			}
			fmt.Println("Unesite velicinu stranice: ")
			r := bufio.NewReader(os.Stdin)
			unos, _ := r.ReadString('\n')
			unos = strings.Replace(unos, "\n", "", 1)
			unos = strings.Replace(unos, "\r", "", 1)
			velicina, err := strconv.ParseInt(unos, 10, 0)
			if err != nil {
				fmt.Println("Niste uneli broj!")
				break
			}
			fmt.Println("Unesite redni broj stranice koju zelite: ")
			r = bufio.NewReader(os.Stdin)
			unos, _ = r.ReadString('\n')
			unos = strings.Replace(unos, "\n", "", 1)
			unos = strings.Replace(unos, "\r", "", 1)
			redni_broj, err := strconv.ParseInt(unos, 10, 0)
			if err != nil {
				fmt.Println("Niste uneli broj!")
				break
			}
			trazeni_kljucevi := RangeScan.DoRangeScan(kljuc1, kljuc2, int(velicina), int(redni_broj), kljucevi_memtable)
			vrednosti := make([][]byte, len(trazeni_kljucevi))
			for i := 0; i < len(trazeni_kljucevi); i++ {
				if engine.da_li_je_skip {
					_, value := citanje.CitajSkip(trazeni_kljucevi[i], engine.mems, engine.cache)
					vrednosti[i] = value
				} else {
					_, value := citanje.CitajBTree(trazeni_kljucevi[i], engine.memb, engine.cache)
					vrednosti[i] = value
				}
			}
			fmt.Println("Vrednosti dobijene range scan-om su: ")
			fmt.Println(vrednosti)
			break
		case "6":
			desetPlusMeni(engine)
			break
		case "x":
			if engine.da_li_je_skip {
				b, velicina := engine.mems.ProveriFlush()
				if b == false && velicina > 0 {
					engine.mems.Flush(engine.indexBloom, engine.konfiguracije["da_li_je_vise_fajlova"])
					engine.indexBloom++
				}
				novi_pod := strconv.Itoa(engine.levelBloom) + "\n" + strconv.Itoa(engine.indexBloom)
				file, err := os.OpenFile("Data/Konfiguracije/podaci.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					panic(err)
				}
				file.WriteString(novi_pod)
				file.Close()

			} else {
				b, velicina := engine.memb.ProveriFlush()
				if b == false && velicina > 0 {
					engine.memb.Flush(engine.indexBloom, engine.konfiguracije["da_li_je_vise_fajlova"])
					engine.indexBloom++
				}
				novi_pod := strconv.Itoa(engine.levelBloom) + "\n" + strconv.Itoa(engine.indexBloom)
				file, err := os.OpenFile("Data/Konfiguracije/podaci.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					panic(err)
				}
				file.WriteString(novi_pod)
				file.Close()
			}
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
		fmt.Println("6:PROCENA IZ HLL")
		fmt.Println("7:NAPRAVI SIM HASH")
		fmt.Println("8:UPOREDI DVA SIM HASH")
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
			addCms(engine)
			break
		case "3":
			fmt.Println(checkCms(engine))
			break
		case "4":
			//key := nabavi_vrednosti_brisanje()
			makeHll(engine)
			break
		case "5":
			if engine.hll == nil {
				fmt.Println("Morate prvo napraviti hll!")
				break
			}
			addHll(engine)
			break
		case "6":
			estimateHll(engine)
			break
		case "7":
			makeSimHash(engine)
			break
		case "8":
			SimHashDistance(engine)
			break
		case "x":
			b = false
			break
		default:
			fmt.Println("Pogresan unos")
		}
	}
}

func makeSimHash(engine *Engine) {
	key := nabavi_vrednosti_brisanje()
	fmt.Println("Unesite tekst: ")
	r := bufio.NewReader(os.Stdin)
	tekst, _ := r.ReadString('\n')
	tekst = strings.Replace(tekst, "\n", "", 1)
	tekst = strings.Replace(tekst, "\r", "", 1)
	vrednost := SimHash.SimHash(tekst)
	engine.simHashevi[key] = vrednost
	engine.simHashKljucevi = append(engine.simHashKljucevi, key)
}

func SimHashDistance(engine *Engine) {
	fmt.Println("Prvi sim hash:")
	kljuc1 := nabavi_vrednosti_brisanje()
	postoji := false
	for i := 0; i < len(engine.simHashKljucevi); i++ {
		if kljuc1 == engine.simHashKljucevi[i] {
			postoji = true
		}
	}
	if !postoji {
		fmt.Println("Sim hash pod datim kljucem ne postoji")
		return
	}
	fmt.Println("Drugi sim hash:")
	kljuc2 := nabavi_vrednosti_brisanje()
	postoji2 := false
	for i := 0; i < len(engine.simHashKljucevi); i++ {
		if kljuc2 == engine.simHashKljucevi[i] {
			postoji2 = true
		}
	}
	if !postoji2 {
		fmt.Println("Sim hash pod datim kljucem ne postoji")
		return
	}
	fmt.Println(SimHash.Hamming(engine.simHashevi[kljuc1], engine.simHashevi[kljuc2]))
}

func makeCms(engine *Engine) {
	key := nabavi_vrednosti_brisanje()
	cms := CountMinSketch.CreateCMS(0.1, 0.1)
	engine.cms_pokaz[key] = cms
	engine.cms_podaci[key] = cms.Bytes
	engine.cms = append(engine.cms, key)
}

func addCms(engine *Engine) {
	kljuc_cms := nabavi_vrednosti_brisanje()
	postoji := false
	for i := 0; i < len(engine.cms); i++ {
		if kljuc_cms == engine.cms[i] {
			postoji = true
		}
	}
	if !postoji {
		fmt.Println("CMS pod datim kljucem ne postoji")
		return
	}
	key := nabavi_vrednosti_brisanje()
	bajtovi1 := engine.cms_pokaz[kljuc_cms].Bytes
	podaci := CountMinSketch.Add(key, engine.cms_pokaz[kljuc_cms].Hashes, int(engine.cms_pokaz[kljuc_cms].M), bajtovi1)
	engine.cms_podaci[kljuc_cms] = podaci
	engine.cms_pokaz[kljuc_cms].Bytes = podaci
	bajtovi := make([]byte, 0)
	for i := 0; i < len(engine.cms); i++ {
		kljucic := engine.cms[i]
		bajtovi = append(bajtovi, engine.cms_podaci[kljucic]...)
	}

	//NE RADI UPIS U FAJL NA NEKIM LAPTOPOVIMA!!!!!!!!!!!!!!!

	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	file_cms, errData := os.OpenFile(path+"/cms.txt", os.O_CREATE|os.O_WRONLY, 0777)
	if errData != nil {
		panic(errData)
	}
	file_cms.Write(bajtovi)
	// os.WriteFile("Strukture//CountMinSketch//cms.txt", bajtovi,os.O_RDWR
	// 0777)
}

func checkCms(engine *Engine) int {
	kljuc_cms := nabavi_vrednosti_brisanje()
	postoji := false
	for i := 0; i < len(engine.cms); i++ {
		if kljuc_cms == engine.cms[i] {
			postoji = true
		}
	}
	if !postoji {
		fmt.Println("CMS pod datim kljucem ne postoji")
		return -5
	}
	key := nabavi_vrednosti_brisanje()
	return CountMinSketch.Cms(key, engine.cms_pokaz[kljuc_cms].Hashes, int(engine.cms_pokaz[kljuc_cms].M), engine.cms_podaci[kljuc_cms])
}

func addHll(engine *Engine) {
	kljuc_hll := nabavi_vrednosti_brisanje()
	postoji := false
	for i := 0; i < len(engine.hll); i++ {
		if kljuc_hll == engine.hll[i] {
			postoji = true
		}
	}
	if !postoji {
		fmt.Println("HLL pod datim kljucem ne postoji")
		return
	}
	key := nabavi_vrednosti_brisanje()
	podaci := HyperLogLog.Add(key, engine.hll_pokaz[kljuc_hll].Reg, engine.hll_pokaz[kljuc_hll].P)
	engine.hll_podaci[kljuc_hll] = podaci
	engine.hll_pokaz[kljuc_hll].Reg = podaci
	bajtovi := make([]byte, 0)
	for i := 0; i < len(engine.hll); i++ {
		kljucic := engine.hll[i]
		bajtovi = append(bajtovi, engine.hll_podaci[kljucic]...)
	}

	//NE RADI UPIS U FAJL NA NEKIM LAPTOPOVIMA!!!!!!!!!!!!!!!

	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	file_hll, errData := os.OpenFile(path+"/hll.txt", os.O_CREATE|os.O_WRONLY, 0777)
	if errData != nil {
		panic(errData)
	}
	file_hll.Write(bajtovi)
	fmt.Println("Element je uspesno dodat u hll!")
}

func estimateHll(engine *Engine) {
	fmt.Println("Procena broja elemenata u hll: ")
	kljuc_hll := nabavi_vrednosti_brisanje()
	postoji := false
	for i := 0; i < len(engine.hll); i++ {
		if kljuc_hll == engine.hll[i] {
			postoji = true
		}
	}
	if !postoji {
		fmt.Println("HLL pod datim kljucem ne postoji")
	}
	fmt.Println(engine.hll_pokaz[kljuc_hll].Estimate())
}

func makeHll(engine *Engine) {
	kljuc := nabavi_vrednosti_brisanje()
	hll := HyperLogLog.MakeHyperLogLog(HyperLogLog.HLL_MAX_PRECISION)
	engine.hll_pokaz[kljuc] = hll
	engine.cms_podaci[kljuc] = hll.Reg
	engine.hll = append(engine.hll, kljuc)
}
