package RangeScan

import (
	citanje "Operacije/Citanje"
	"Strukture/SSTable"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func DoRangeScan(min string, maks string, velicina int, broj int, kljucevi_memtable []string, da_li_je_skip bool) [][]byte {
	path1, _ := filepath.Abs("../Spojeno/Data/SSTableData")
	path := strings.ReplaceAll(path1, `\`, "/")
	_, _, index_files, _, _ := citanje.Svi_fajlovi(path)
	var kljucevi1 []string
	for i := range index_files {
		indFile, err := os.OpenFile(path+"/"+index_files[i], os.O_RDONLY, 0666)
		if err != nil {
			panic(err)
		}
		kljucevi1 = append(kljucevi1, SSTable.Svi_kljucevi_jednog_fajla(indFile)...)
	}
	kljucevi1 = append(kljucevi1, kljucevi_memtable...)
	kljucevi := make([]string, 0)
	sort.Strings(kljucevi1)
	potrebni_kljucevi := make([]string, 0)
	indeks := velicina * (broj - 1)
	for i := indeks; i < indeks+velicina; i++ {
		potrebni_kljucevi = append(potrebni_kljucevi, kljucevi[i])
	}
	//trazene_vrednosti := make([][]byte, len(potrebni_kljucevi))
	for i := 0; i < len(potrebni_kljucevi); i++ {
		//trazene_vrednosti = append(trazene_vrednosti, )
	}
	return make([][]byte, 0)
}
