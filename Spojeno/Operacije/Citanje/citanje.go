package citanje

import (
	"Spojeno/Strukture/MemTableBTree"
	"Strukture/Cache"
	"Strukture/MemTableSkipList"
	"Strukture/SSTable"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CitajSKip(kljuc string, memTable *MemTableSkipList.MemTable, cache *Cache.Cache) (bool, []byte) {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	data_files, _, index_files, summary_files, _ := Svi_fajlovi(path)
	b, value := memTable.NadjiElement(kljuc)
	if b {
		return b, value
	} else {
		b, _ := cache.GetFromCache(kljuc) //ako mu pristupi stavi ga na pocetak cache-a
		if b {
			value, _ = cache.NadjiUCache(kljuc)
			return b, value
		} else {
			fmt.Println("Usao u summary fajlove")
			// for _, bfajl := range filter_files {
			// 	BloomFilter
			// }
			sumBr := 0
			for _, sFajl := range summary_files {
				sumBr++
				sumFile, err := os.OpenFile(path+"/SSTableData/"+sFajl, os.O_RDONLY, 0777)
				if err != nil {
					panic(err)
				}
				offset, b := SSTable.NadjiSummary(kljuc, sumFile)
				if b {
					fmt.Println("Nasao u summary fajlu")
					indBr := 0
					for _, iFajl := range index_files {
						indBr++
						if indBr == sumBr {
							fmt.Println("Cita u indexu")
							indFile, err := os.OpenFile(path+"/SSTableData/"+iFajl, os.O_RDONLY, 0777)
							if err != nil {
								panic(err)
							}
							b, offset1 := SSTable.NadjiIndex(offset, indFile, kljuc)
							if b {
								fmt.Println("Nasao u indexu")
								datBr := 0
								for _, dataFajl := range data_files {
									datBr++
									if datBr == indBr {
										fmt.Println("Cita u data")
										dataFile, err := os.OpenFile(path+"/SSTableData/"+dataFajl, os.O_RDONLY, 0777)
										if err != nil {
											panic(err)
										}
										b, value := SSTable.NadjiElement(offset1, dataFile, kljuc)
										if b {
											fmt.Println("Nasao u data")
											return b, value

										}
									}
								}
							}
						}

					}

				}
			}
		}

	}
	return false, nil
}
func CitajBTree(kljuc string, memTable *MemTableBTree.MemTable, cache *Cache.Cache) (bool, []byte) {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	data_files, _, index_files, summary_files, _ := Svi_fajlovi(path)
	b, value := memTable.NadjiElement(kljuc)
	if b {
		return b, value
	} else {
		b, _ := cache.GetFromCache(kljuc) //ako mu pristupi stavi ga na pocetak cache-a
		if b {
			value, _ = cache.NadjiUCache(kljuc)
			return b, value
		} else {
			fmt.Println("Usao u summary fajlove")
			// for _, bfajl := range filter_files {
			// 	BloomFilter
			// }
			sumBr := 0
			for _, sFajl := range summary_files {
				sumBr++
				sumFile, err := os.OpenFile(path+"/SSTableData/"+sFajl, os.O_RDONLY, 0777)
				if err != nil {
					panic(err)
				}
				offset, b := SSTable.NadjiSummary(kljuc, sumFile)
				if b {
					fmt.Println("Nasao u summary fajlu")
					indBr := 0
					for _, iFajl := range index_files {
						indBr++
						if indBr == sumBr {
							fmt.Println("Cita u indexu")
							indFile, err := os.OpenFile(path+"/SSTableData/"+iFajl, os.O_RDONLY, 0777)
							if err != nil {
								panic(err)
							}
							b, offset1 := SSTable.NadjiIndex(offset, indFile, kljuc)
							if b {
								fmt.Println("Nasao u indexu")
								datBr := 0
								for _, dataFajl := range data_files {
									datBr++
									if datBr == indBr {
										fmt.Println("Cita u data")
										dataFile, err := os.OpenFile(path+"/SSTableData/"+dataFajl, os.O_RDONLY, 0777)
										if err != nil {
											panic(err)
										}
										b, value := SSTable.NadjiElement(offset1, dataFile, kljuc)
										if b {
											fmt.Println("Nasao u data")
											return b, value

										}
									}
								}
							}
						}

					}

				}
			}
		}

	}
	return false, nil
}

func Svi_fajlovi(folder string) ([]string, []string, []string, []string, []string) {
	fajlovi, err := ioutil.ReadDir(folder)
	if err != nil {
		panic(err)
	}
	var svi []string
	for _, file := range fajlovi {
		if file.IsDir() {
			fajlovi1, err := ioutil.ReadDir(folder + "/" + file.Name())
			if err != nil {
				panic(err)
			}
			for _, file1 := range fajlovi1 {
				svi = append(svi, file1.Name())
			}
		} else {
			svi = append(svi, file.Name())
		}
	}
	var data_files []string
	var filter_files []string
	var index_files []string
	var summary_files []string
	var toc_files []string
	for f := range svi {
		if strings.Contains(svi[f], "DataFileL") {
			data_files = append(data_files, svi[f])
		}
		if strings.Contains(svi[f], "FilterFileL") {
			filter_files = append(filter_files, svi[f])
		}
		if strings.Contains(svi[f], "IndexFileL") {
			index_files = append(index_files, svi[f])
		}
		if strings.Contains(svi[f], "SummaryFileL") {
			summary_files = append(summary_files, svi[f])
		}
		if strings.Contains(svi[f], "TocFileL") {
			toc_files = append(toc_files, svi[f])
		}

	}
	return data_files, filter_files, index_files, summary_files, toc_files
}
