// package citanje

// import (
// 	"io/ioutil"
// 	"path/filepath"
// 	"strings"
// )

// func citaj(key string, memTable *MemTableSkipList) {
// 	path1, _ := filepath.Abs("../Key-Value-engine/Data")
// 	path := strings.ReplaceAll(path1, `\`, "/")
// 	data_files, filter_files, index_files, summary_files, _ := Svi_fajlovi(path)

// 	for elem := range filter_files {
// 		print(elem)
// 	}

// }
// func Svi_fajlovi(folder string) ([]string, []string, []string, []string, []string) {
// 	fajlovi, err := ioutil.ReadDir(folder)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var svi []string
// 	for _, file := range fajlovi {
// 		if file.IsDir() {
// 			fajlovi1, err := ioutil.ReadDir(folder + "/" + file.Name())
// 			if err != nil {
// 				panic(err)
// 			}
// 			for _, file1 := range fajlovi1 {
// 				svi = append(svi, file1.Name())
// 			}
// 		} else {
// 			svi = append(svi, file.Name())
// 		}
// 	}
// 	var data_files []string
// 	var filter_files []string
// 	var index_files []string
// 	var summary_files []string
// 	var toc_files []string
// 	for f := range svi {
// 		if strings.Contains(svi[f], "DataFileL") {
// 			data_files = append(data_files, svi[f])
// 		}
// 		if strings.Contains(svi[f], "FilterFileL") {
// 			filter_files = append(filter_files, svi[f])
// 		}
// 		if strings.Contains(svi[f], "IndexFileL") {
// 			index_files = append(index_files, svi[f])
// 		}
// 		if strings.Contains(svi[f], "SummaryFileL") {
// 			summary_files = append(summary_files, svi[f])
// 		}
// 		if strings.Contains(svi[f], "TocFileL") {
// 			toc_files = append(toc_files, svi[f])
// 		}

// 	}
// 	return data_files, filter_files, index_files, summary_files, toc_files
// }
