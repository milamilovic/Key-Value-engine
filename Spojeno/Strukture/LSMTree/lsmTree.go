package LSMTree

import (
	"Strukture/SSTable"
	"io/ioutil"

	//"path/filepath"
	"strconv"
	"strings"
)

//func main() {
//	path1, _ := filepath.Abs("../Key-Value-engine/Data")
//	path := strings.ReplaceAll(path1, `\`, "/")
//	stablo := Kreiraj_lsmTree(4, 1, 8)
//	kompakcije(*stablo, 1, path)
//}

type Lsm struct {
	najveci_nivo     int
	najveca_velicina int
	najveci_listLen  int
}

func Kreiraj_lsmTree(nivo int, velicina int, listLen int) *Lsm {
	LSM := Lsm{nivo, velicina, listLen}
	return &LSM
}

func Svi_fajlovi(n int, folder string) ([]string, []string, []string, []string, []string) {
	nivo := strconv.Itoa(n)
	fajlovi, err := ioutil.ReadDir(folder)
	if err != nil {
		panic(err)
	}
	//for _, file := range fajlovi {
	//	fmt.Println(file.Name(), file.IsDir())
	var svi []string
	for _, file := range fajlovi {
		//fmt.Println(file.Name(), file.IsDir())
		if file.IsDir() {
			fajlovi1, err := ioutil.ReadDir(folder + "/" + file.Name())
			if err != nil {
				panic(err)
			}
			for _, file1 := range fajlovi1 {
				//fmt.Println(file1.Name())
				svi = append(svi, file1.Name())
			}
		} else {
			//fmt.Println(file.Name())
			svi = append(svi, file.Name())
		}
	}
	//for f := range svi {
	//	fmt.Println(svi[f])
	//}
	var data_files []string
	var filter_files []string
	var index_files []string
	var summary_files []string
	var toc_files []string
	for f := range svi {
		if strings.Contains(svi[f], "DataFileL"+nivo+"Id") {
			data_files = append(data_files, svi[f])
		}
		if strings.Contains(svi[f], "FilterFileL"+nivo+"Id") {
			filter_files = append(filter_files, svi[f])
		}
		if strings.Contains(svi[f], "IndexFileL"+nivo+"Id") {
			index_files = append(index_files, svi[f])
		}
		if strings.Contains(svi[f], "SummaryFileL"+nivo+"Id") {
			summary_files = append(summary_files, svi[f])
		}
		if strings.Contains(svi[f], "TocFileL"+nivo+"Id") {
			toc_files = append(toc_files, svi[f])
		}

	}
	return data_files, filter_files, index_files, summary_files, toc_files
}

func Da_li_nastavljamo(LSM Lsm, n int, folder string) ([]string, []string, []string, []string, []string, bool) {
	data_files, filter_files, index_files, summary_files, toc_files := Svi_fajlovi(n, folder)
	var da_ne bool
	if len(data_files) == LSM.najveca_velicina && len(data_files) > 0 {
		da_ne = false //znaci da ne nastavljamo kompakcije
	} else {
		da_ne = true //nastavice se sa kompakcijama
	}
	return data_files, filter_files, index_files, summary_files, toc_files, da_ne
}

func Kompakcije(LSM Lsm, n int, folder string) {
	if n >= LSM.najveci_nivo {
		return // dosli smo do poslednjeg nivoa
	}
	//data_files, filter_files, index_files, summary_files, toc_files, da_ne := da_li_nastavljamo(LSM, n, folder)
	data_files, _, _, _, _, da_ne := Da_li_nastavljamo(LSM, n, folder)
	if da_ne == false {
		return // ne radimo kompakciju
	}
	// else , radimo kompakciju
	SSTable.Kompakcija(len(data_files), LSM.najveci_nivo, n, LSM.najveci_listLen)
}
