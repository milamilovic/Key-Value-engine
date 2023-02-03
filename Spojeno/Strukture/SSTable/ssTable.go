package SSTable

import (
	"Strukture/BloomFilter"
	"Strukture/SkipList"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func NapraviSSTable(lCvor []*SkipList.SkipListNode, level int, index int) {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
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

	var offsetInd uint64 = 0

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
	sumFile.Close()

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

		size := 4 + 8 + 1 + 8 + 8 + key_u + val_u
		offset_ind := make([]byte, 8)
		binary.LittleEndian.PutUint64(offset_ind, offsetInd)

		indFile.Write(keySize)
		indFile.Write([]byte(cvor.GetKey()))
		indFile.Write(offset_ind)

		offsetInd = offsetInd + size

	}
	datFile.Close()
	indFile.Close()
}

func NadjiSummary(kljuc string, f *os.File) bool {
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
		return false
	}
	return true

}
func NadjiIndex(f *os.File, kljuc string) (bool, uint64) {
	f.Seek(0, 0)
	for true {
		size := make([]byte, 8)
		f.Read(size)
		keySize := binary.LittleEndian.Uint64(size) //dobijamo velicinu kljuca
		key_read := make([]byte, keySize)
		f.Read(key_read)
		key := string(key_read)
		fmt.Println(key)
		of := make([]byte, 8)
		f.Read(of)
		offset := binary.LittleEndian.Uint64(of) //velicina offseta
		if key == kljuc {
			return true, offset
		}
		if key > kljuc {
			return false, 0
		}
		if key < kljuc {
			continue
		}
	}
	return false, 0
}

func NadjiElement(offset uint64, f *os.File, kljuc string) (bool, []byte) {
	f.Seek(int64(offset), 0)

	bytes := make([]byte, 8)
	f.Read(bytes)
	bytes = make([]byte, 4)
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

func Kompakcija(brojFajlova int, maxLevel int, level int, listLen int) {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	br := 0
	for br < brojFajlova {

		bloom := BloomFilter.New_bloom(listLen, 0.05, level+1, br)
		l := strconv.Itoa(bloom.Level)
		i := strconv.Itoa(bloom.Index)
		skipList := SkipList.NapraviSkipList(listLen) //velicina koja
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
			_, time1, tomb1, key_s1, val_s1, key1, val1, end1 := GetData(f1)
			_, time2, tomb2, key_s2, val_s2, key2, val2, end2 := GetData(f2)
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
						BloomFilter.Add(key2, path+"/SSTableData/filterFileL"+l+"Id"+i+".txt")

					} else if time1 > time2 {
						skipList.Add(key1, val1)
						BloomFilter.Add(key1, path+"/SSTableData/filterFileL"+l+"Id"+i+".txt")

					} else {
						skipList.Add(key1, val1)
						BloomFilter.Add(key1, path+"/SSTableData/filterFileL"+l+"Id"+i+".txt")

					}
				}
			} else if key1 > key2 {
				if tomb2[0] != 1 {
					skipList.Add(key2, val2)
					BloomFilter.Add(key2, path+"/SSTableData/filterFileL"+l+"Id"+i+".txt")

				}
				f1.Seek(int64(4+8+1+8+8+key_s1+val_s1), 1)

			} else {
				if tomb1[0] != 1 {
					skipList.Add(key1, val1)
					BloomFilter.Add(key1, path+"/SSTableData/filterFileL"+l+"Id"+i+".txt")

				}
				f2.Seek(int64(4+8+1+8+8+key_s2+val_s2), 1)
			}
		}
		for {
			_, _, _, key_s1, val_s1, key1, val1, end1 := GetData(f1)
			if end1 == true {
				break
			}
			skipList.Add(key1, val1)
			f1.Seek(int64(4+8+1+8+8+key_s1+val_s1), 1)
		}

		for {
			_, _, _, key_s2, val_s2, key2, val2, end2 := GetData(f2)
			if end2 == true {
				break
			}
			skipList.Add(key2, val2)
			f1.Seek(int64(4+8+1+8+8+key_s2+val_s2), 1)
		}
		index := NovoImeFajla(level + 1)
		f1.Close()
		f2.Close()
		newData := skipList.GetElements()
		NapraviSSTable(newData, level+1, index)
		ObrisiFajlove(level, br)
		NapraviTOC(level+1, index)
	}
	index := NovoImeFajla(level + 1)
	if index == brojFajlova {
		if (level + 1) < maxLevel {
			Kompakcija(brojFajlova, maxLevel, level+1, 2*listLen)
		}
	}
}

func NapraviTOC(level int, index int) {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	tocFile, err := os.OpenFile(path+"/TOCFiles/TocFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index-1)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	_, er := tocFile.Write([]byte(path + "/SSTableData/DataFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db"))
	if er != nil {
		panic(er)
	}
	_, er = tocFile.Write([]byte(path + "/SSTableData/IndexFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db"))
	if er != nil {
		panic(er)
	}
	_, er = tocFile.Write([]byte(path + "/SSTableData/SummaryFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db"))
	if er != nil {
		panic(er)
	}
	_, er = tocFile.Write([]byte(path + "/SSTableData/FilterFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db"))
	if er != nil {
		panic(er)
	}
	_, er = tocFile.Write([]byte(path + "/SSTableData/MerkleL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".txt"))
	if er != nil {
		panic(er)
	}

}

func ObrisiFajlove(level int, index int) {
	path1, _ := filepath.Abs("../Spojeno/Data")
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
	err7 := os.Remove(path + "/SSTableData/FilterFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".db")
	if err7 != nil {
		panic(err7)
	}
	err8 := os.Remove(path + "/SSTableData/FilterFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db")
	if err8 != nil {
		panic(err8)
	}
	err9 := os.Remove(path + "/SSTableData/MerkleL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".txt")
	if err9 != nil {
		panic(err7)
	}
	err10 := os.Remove(path + "/SSTableData/MerkleL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".txt")
	if err10 != nil {
		panic(err8)
	}

}

func NovoImeFajla(level int) int {
	path1, _ := filepath.Abs("../Spojeno/Data")
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

func GetData(f *os.File) (uint32, uint64, []byte, uint64, uint64, string, []byte, bool) {

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

func Svi_kljucevi_jednog_fajla(f *os.File) []string {
	size := make([]byte, 8)
	f.Read(size)
	keySize := binary.LittleEndian.Uint64(size) //dobijamo velicinu kljuca
	key_read := make([]byte, keySize)
	f.Read(key_read)
	//key1 := string(key_read) //prvi kljuc
	size = make([]byte, 8)
	f.Read(size)
	keySize = binary.LittleEndian.Uint64(size) //dobijamo velicinu kljuca
	key_read = make([]byte, keySize)
	f.Read(key_read)
	//key2 := string(key_read) //poslednji kljuc
	var kljucevi []string
	for true { //citamo kljuceve redom
		size := make([]byte, 8)
		_, err := f.Read(size)

		if err != nil {
			break
		}
		keySize := binary.LittleEndian.Uint64(size)
		fmt.Println(keySize)
		key_read := make([]byte, keySize)
		_, err = f.Read(key_read)
		fmt.Println(key_read)
		if err != nil {
			break
		}
		key := string(key_read)
		kljucevi = append(kljucevi, key)
	}
	return kljucevi

}
