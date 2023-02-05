package RangeScan

import (
	citanje "Operacije/Citanje"
	"Strukture/SSTable"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func DoRangeScan(min string, maks string, velicina int, broj int, kljucevi_memtable []string) []string {
	path1, _ := filepath.Abs("../Spojeno/Data/SSTableData")
	path := strings.ReplaceAll(path1, `\`, "/")
	_, _, index_files, _, _ := citanje.Svi_fajlovi(path)
	var kljucevi1 []string
	for i := range index_files {
		indFile, err := os.OpenFile(path+"/"+index_files[i], os.O_RDONLY, 0666)
		if err != nil {
			path1, _ = filepath.Abs("../Projekat/Spojeno/Data/SSTableData")
			path = strings.ReplaceAll(path1, `\`, "/")
			indFile, err = os.OpenFile(path+"/"+index_files[i], os.O_RDONLY, 0666)
			if err != nil {
				panic(err)
			}
		}
		kljucevi1 = append(kljucevi1, SSTable.Svi_kljucevi_jednog_fajla(indFile)...)
	}
	kljucevi1 = append(kljucevi1, kljucevi_memtable...)
	sort.Strings(kljucevi1)
	kljucevi := make([]string, 0)
	for i := 0; i < len(kljucevi1); i++ {
		if kljucevi1[i] < maks && kljucevi1[i] > min {
			kljucevi = append(kljucevi, kljucevi1[i])
		}
	}
	potrebni_kljucevi := make([]string, 0)
	indeks := velicina * (broj - 1)
	for i := indeks; i < indeks+velicina; i++ {
		if i < len(kljucevi) {
			potrebni_kljucevi = append(potrebni_kljucevi, kljucevi[i])
		}
	}
	return potrebni_kljucevi
}

func DoRangeScanJedanFajl(min string, maks string, velicina int, broj int, kljucevi_memtable []string) []string {
	path1, _ := filepath.Abs("../Spojeno/Data/SSTableData")
	path := strings.ReplaceAll(path1, `\`, "/")
	data_files, _, _, _, _ := citanje.Svi_fajlovi(path)
	var kljucevi1 []string
	for i := range data_files {
		indFile, err := os.OpenFile(path+"/"+data_files[i], os.O_RDONLY, 0666)
		if err != nil {
			path1, _ = filepath.Abs("../Projekat/Spojeno/Data/SSTableData")
			path = strings.ReplaceAll(path1, `\`, "/")
			indFile, err = os.OpenFile(path+"/"+data_files[i], os.O_RDONLY, 0666)
			if err != nil {
				panic(err)
			}
		}
		kljucevi_iz_fajla := SSTable.Svi_kljucevi_jednog_fajla_jedan_fajl(indFile)
		kljucevi1 = append(kljucevi1, kljucevi_iz_fajla...)
	}
	kljucevi1 = append(kljucevi1, kljucevi_memtable...)
	sort.Strings(kljucevi1)
	kljucevi := make([]string, 0)
	for i := 0; i < len(kljucevi1); i++ {
		if kljucevi1[i] < maks && kljucevi1[i] > min {
			kljucevi = append(kljucevi, kljucevi1[i])
		}
	}
	potrebni_kljucevi := make([]string, 0)
	indeks := velicina * (broj - 1)
	for i := indeks; i < indeks+velicina; i++ {
		if i < len(kljucevi) {
			potrebni_kljucevi = append(potrebni_kljucevi, kljucevi[i])
		}
	}
	return potrebni_kljucevi
}
