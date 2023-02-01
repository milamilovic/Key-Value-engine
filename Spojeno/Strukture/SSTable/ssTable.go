package SSTable

import (
	"Strukture/BloomFilter"
	"Strukture/MerkleTree"
	"Strukture/SkipList"
	"encoding/binary"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func MakeSSTable(lCvor []*SkipList.SkipListNode, level int, index int) {
	path1, _ := filepath.Abs("../Key-Value-engine/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	//fmt.Println(path)
	datFile, errData := os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if errData != nil {
		panic(errData)
	}
	indFile, errInd := os.OpenFile(path+"/SSTableData/IndexFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if errInd != nil {
		panic(errInd)
	}
	sumFile, errSum := os.OpenFile(path+"/SSTableData/SummaryFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if errSum != nil {
		panic(errInd)
	}
	offsetDat := 0
	offsetInd := 0

	first := make([]byte, 8)
	first_u := uint64(len(lCvor[0].GetKey()))
	binary.LittleEndian.PutUint64(first, first_u)

	last := make([]byte, 8)
	last_u := uint64(len(lCvor[len(lCvor)-1].GetKey()))
	binary.LittleEndian.PutUint64(last, last_u)

	sumFile.Write(first)
	sumFile.Write([]byte(lCvor[0].GetKey()))
	sumFile.Write(last)
	sumFile.Write([]byte(lCvor[len(lCvor)-1].GetKey()))

	for _, cvor := range lCvor {
		crc := make([]byte, 4)
		binary.BigEndian.PutUint32(crc, uint32(crc32.ChecksumIEEE(cvor.GetValue())))

		timestamp := make([]byte, 8)
		binary.BigEndian.PutUint64(timestamp, uint64(cvor.GetTimeStamp()))

		tombstone := make([]byte, 1)
		if cvor.GetTombstone() {
			tombstone[0] = 1
		}
		keySize := make([]byte, 8)
		key_u := uint64(len(cvor.GetKey()))
		binary.LittleEndian.PutUint64(keySize, key_u)

		valSize := make([]byte, 8)
		val_u := uint64(len(cvor.GetValue()))
		binary.LittleEndian.PutUint64(valSize, val_u)

		datFile.Write(crc)
		datFile.Write(timestamp)
		datFile.Write(tombstone)
		datFile.Write(keySize)
		datFile.Write(valSize)
		datFile.Write([]byte(cvor.GetKey()))
		datFile.Write(cvor.GetValue())

		size := 4 + 8 + 1 + 8 + 8 + len(cvor.GetKey()) + len(cvor.GetValue())
		offset_ind := make([]byte, 8)
		binary.LittleEndian.PutUint64(offset_ind, uint64(offsetDat))
		indFile.Write(keySize)
		indFile.Write([]byte(cvor.GetKey()))
		indFile.Write(offset_ind)
		offsetDat = offsetDat + int(size)
		ind_offset := 16 + key_u

		offset_sum := make([]byte, 8)
		binary.LittleEndian.PutUint64(offset_sum, uint64(offsetInd))
		sumFile.Write(keySize)
		sumFile.Write([]byte(cvor.GetKey()))
		sumFile.Write(offset_sum)
		offsetInd = offsetInd + int(ind_offset)
	}

	datFile.Close()
	indFile.Close()
	sumFile.Close()
}

func nadjiSummary(kljuc string, f *os.File) (uint64, bool) {
	size := make([]byte, 8)
	f.Read(size)
	keySize := binary.LittleEndian.Uint64(size) //dobijamo velicinu kljuca
	key_read := make([]byte, keySize)
	f.Read(key_read)
	key1 := string(key_read) //prvi kljuc
	size = make([]byte, 8)
	f.Read(size)
	keySize = binary.LittleEndian.Uint64(size) //dobijamo velicinu kljuca
	key_read = make([]byte, keySize)
	f.Read(key_read)
	key2 := string(key_read) //poslednji kljuc
	if kljuc > key2 || kljuc < key1 {
		return 0, false
	}

	for true { //citamo kljuceve redom
		size = make([]byte, 8)
		f.Read(size)
		keySize = binary.LittleEndian.Uint64(size)
		key_read = make([]byte, keySize)
		f.Read(key_read)
		key := string(key_read)
		if key > kljuc {
			return 0, false
		} else if key < kljuc {
			continue
		} else { //isti su, nasli smo ga
			offset := make([]byte, 8) //citamo sledecih 8 bajtova sto je nas trazeni offset
			f.Read(offset)
			offsetSize := binary.LittleEndian.Uint64(size)
			return offsetSize, true
		}
	}
	return 0, false

}
func nadjiIndex(offset uint64, f *os.File, kljuc string) (bool, uint64) {

	f.Seek(int64(offset), 0)
	size := make([]byte, 8)
	f.Read(size)
	for true {
		size := make([]byte, 8)
		f.Read(size)
		keySize := binary.LittleEndian.Uint64(size) //dobijamo velicinu kljuca
		key_read := make([]byte, keySize)
		f.Read(key_read)
		key := string(key_read)
		if key > kljuc {
			return false, 0
		} else if key < kljuc {
			continue
		} else { //isti su, nasli smo ga
			offset := make([]byte, 8) //citamo sledecih 8 bajtova sto je nas trazeni offset
			f.Read(offset)
			offsetSize := binary.LittleEndian.Uint64(size)
			return true, offsetSize
		}
	}
	return false, 0
}

func nadjiElement(offset uint64, f *os.File, kljuc string) (bool, []byte) {
	f.Seek(int64(offset), 0)
	bytes := make([]byte, 8)
	f.Read(bytes)
	t := make([]byte, 1)
	f.Read(t)
	if t[0] == 1 {
		return false, nil
	}
	size := make([]byte, 8)
	f.Read(size)
	keySize := binary.LittleEndian.Uint64(size)
	size = make([]byte, 8)
	f.Read(size)
	valueSize := binary.LittleEndian.Uint64(size)

	key_read := make([]byte, keySize)
	f.Read(key_read)
	value_read := make([]byte, valueSize)
	f.Read(value_read)
	value := value_read
	return true, value
}

func Compaction(brojFajlova int, maxLevel int, level int, listLen int) {
	br := 0
	path1, _ := filepath.Abs("../Key-Value-engine/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	for br < brojFajlova {
		bloomFilter := BloomFilter.New_bloom(brojFajlova, 2)
		skipList := SkipList.MakeSkipList(10)
		br++
		f1, err := os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(br)+".db", os.O_RDONLY, 0777)
		if err != nil {
			panic(err)
		}
		br++
		f2, err := os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(br)+".db", os.O_RDONLY, 0777)
		if err != nil {
			panic(err)
		}
		for {
			_, time1, tomb1, key_s1, val_s1, key1, val1, end1 := getData(f1)
			_, time2, tomb2, key_s2, val_s2, key2, val2, end2 := getData(f2)
			if end1 == false {
				f1.Seek(-(4 + 8 + 1 + 8 + 8 + int64(key_s1) + int64(val_s1)), 1)
				break
			}
			if end2 == false {
				f1.Seek(-(4 + 8 + 1 + 8 + 8 + int64(key_s2) + int64(val_s2)), 1)
				break
			}

			if key1 == key2 {
				if tomb1[0] == tomb2[0] && tomb1[0] != 1 { //ne upisujemo ako je obrisan
					if time1 < time2 {
						skipList.Add(key2, val2)
						bloomFilter.Add(key2)
					} else if time1 > time2 {
						skipList.Add(key1, val1)
						bloomFilter.Add(key1)
					} else {
						skipList.Add(key1, val1)
						bloomFilter.Add(key1)
					}
				}
			} else if key1 > key2 {
				if tomb2[0] != 1 {
					skipList.Add(key2, val2)
					bloomFilter.Add(key2)
				}
				f1.Seek(int64(4+8+1+8+8+key_s1+val_s1), 1)

			} else {
				if tomb1[0] != 1 {
					skipList.Add(key1, val1)
					bloomFilter.Add(key1)
				}
				f2.Seek(int64(4+8+1+8+8+key_s2+val_s2), 1)
			}
		}
		for {
			_, _, _, key_s1, val_s1, key1, val1, end1 := getData(f1)
			if end1 == true {
				break
			}
			skipList.Add(key1, val1)
			//ovde pukne
			bloomFilter.Add(key1)
			f1.Seek(int64(4+8+1+8+8+key_s1+val_s1), 1)
		}

		for {
			_, _, _, key_s2, val_s2, key2, val2, end2 := getData(f2)
			if end2 == true {
				break
			}
			skipList.Add(key2, val2)
			bloomFilter.Add(key2)
			f1.Seek(int64(4+8+1+8+8+key_s2+val_s2), 1)
		}
		index := newFileName(level+1) + 1
		filterFile, errFil := os.Create(path + "/SSTableData/FilterFileL" + strconv.Itoa(level) +
			"Id" + strconv.Itoa(index-1) + ".db")
		if errFil != nil {
			panic(errFil)

		}
		// bloomFilter.Hashes = nil
		// enc := gob.NewEncoder(filterFile)
		// err = enc.Encode(bloomFilter)
		// if err != nil {
		// 	panic(err)
		// }
		filterFile.Close()
		f1.Close()
		f2.Close()
		newData := skipList.GetElements()
		MakeSSTable(newData, level+1, index)
		deleteFiles(level, br)
		stringovi := []string{}
		for i := 0; i < listLen; i++ {
			stringovi = append(stringovi, newData[i].GetKey())
		}
		podaciZaMerkle := MerkleTree.Pretvori_u_bajtove(stringovi)
		MerkleTree.Kreiraj_MerkleTree(podaciZaMerkle, path+"/SSTableData/MerkleL"+strconv.Itoa(level)+"Id"+strconv.Itoa(index)+".txt")

		writeTOC(level+1, index)
	}
	index := newFileName(level + 1)
	if index == brojFajlova {
		if (level + 1) < maxLevel {
			Compaction(brojFajlova, maxLevel, level+1, 2*listLen)
		}
	}
}

func writeTOC(level int, index int) {
	path1, _ := filepath.Abs("../Key-Value-engine/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	tocFile, err := os.OpenFile(path+"/TOCFiles/TocFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	_, er := tocFile.Write([]byte(path + "/SSTableData/DataFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".db"))
	if er != nil {
		panic(er)
	}
	_, er = tocFile.Write([]byte(path + "/SSTableData/IndexFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".db"))
	if er != nil {
		panic(er)
	}
	_, er = tocFile.Write([]byte(path + "/SSTableData/SummaryFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".db"))
	if er != nil {
		panic(er)
	}
	_, er = tocFile.Write([]byte(path + "/SSTableData/FilterFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".db"))
	if er != nil {
		panic(er)
	}

}

func deleteFiles(level int, index int) {
	path1, _ := filepath.Abs("../Key-Value-engine/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	err1 := os.Remove(path + "/SSTableData/DataFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".db")
	if err1 != nil {
		panic(err1)
	}
	err2 := os.Remove(path + "/SSTableData/DataFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db")
	if err2 != nil {
		panic(err2)
	}
	err3 := os.Remove(path + "/SSTableData/IndexFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".db")
	if err3 != nil {
		panic(err3)
	}
	err4 := os.Remove(path + "/SSTableData/IndexFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db")
	if err4 != nil {
		panic(err4)
	}
	err5 := os.Remove(path + "/SSTableData/SummaryFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".db")
	if err5 != nil {
		panic(err5)
	}
	err6 := os.Remove(path + "/SSTableData/SummaryFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db")
	if err6 != nil {
		panic(err6)
	}
	// err7 := os.Remove("C:/Users/Sonja/Desktop/Key-Value-engine/Data/SSTableData/FilterFileL" + strconv.Itoa(level) +
	// 	"Id" + strconv.Itoa(index) + ".db")
	// if err7 != nil {
	// 	panic(err7)
	// }
	// err8 := os.Remove("C:/Users/Sonja/Desktop/Key-Value-engine/Data/SSTableData/FilterFileL" + strconv.Itoa(level) +
	// 	"Id" + strconv.Itoa(index-1) + ".db")
	// if err8 != nil {
	// 	panic(err8)
	// }

}

func newFileName(level int) int {
	path1, _ := filepath.Abs("../Key-Value-engine/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	br := 1
	for {
		_, err := os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(br)+".db", os.O_RDONLY, 0777)
		if err != nil { //nema takvog fajla, moze da bude novi
			return br
		}
		br++
	}

}

func getData(f *os.File) (uint32, uint64, []byte, uint64, uint64, string, []byte, bool) {

	size := make([]byte, 4)
	_, err := f.Read(size)
	if err == io.EOF {
		return 0, 0, nil, 0, 0, "", nil, true
	}
	crc := binary.LittleEndian.Uint32(size)
	size = make([]byte, 8)
	f.Read(size)
	time := binary.LittleEndian.Uint64(size)
	tomb := make([]byte, 1)
	f.Read(tomb)
	size = make([]byte, 8)
	f.Read(size)
	keySize := binary.LittleEndian.Uint64(size)
	size = make([]byte, 8)
	f.Read(size)
	valueSize := binary.LittleEndian.Uint64(size)
	key := make([]byte, keySize)
	f.Read(key)
	key_s := string(key)
	val := make([]byte, valueSize)
	f.Read(val)

	return crc, time, tomb, keySize, valueSize, key_s, val, false

}
func main() {
	sl := SkipList.MakeSkipList(5)
	sl.Add("1", []byte("a"))
	sl.Add("2", []byte("a"))
	sl.Add("3", []byte("a"))
	sl.Add("4", []byte("a"))
	sl.Add("5", []byte("a"))
	MakeSSTable(sl.GetElements(), 1, 1)

	sl = SkipList.MakeSkipList(5)
	sl.Add("3", []byte("a"))
	sl.Add("7", []byte("a"))
	sl.Add("8", []byte("a"))
	sl.Add("9", []byte("a"))
	sl.Add("5", []byte("a"))
	MakeSSTable(sl.GetElements(), 1, 2)
	writeTOC(1, 1)
	writeTOC(5, 5)
	writeTOC(20, 12)       //ovo radi
	Compaction(2, 2, 1, 5) //bloom filter i merkle brisanje i pravljenje srediti
	// file, err := os.OpenFile("C:/Users/Sonja/Desktop/Key-Value-engine/Data/SSTableData/SummaryFileL1Id1.db", os.O_RDONLY, 0777)
	// if err != nil {
	// 	panic(err)
	// }
	// offset, b := nadjiSummary("2", file)
	// file.Close()
	// if b {
	// 	file, err := os.OpenFile("C:/Users/Sonja/Desktop/Key-Value-engine/Data/SSTableData/IndexFileL1Id1.db", os.O_RDONLY, 0777)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	b, offset = nadjiIndex(offset, file, "2")
	// 	file.Close()
	// 	if b {
	// 		file, err := os.OpenFile("C:/Users/Sonja/Desktop/Key-Value-engine/Data/SSTableData/DataFileL1Id1.db", os.O_RDONLY, 0777)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		bb, _ := nadjiElement(offset, file, "2")
	// 		file.Close()
	// 		if bb {
	// 			print("Nasli smo ga")
	// 		}
	// 	}
	// }

}
