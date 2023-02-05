package SSTable

import (
	"Strukture/BloomFilter"
	"Strukture/MerkleTree"
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
	var datFile *os.File
	datFile, errData := os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if errData != nil {
		path1, _ = filepath.Abs("../Projekat/Spojeno/Data")
		path = strings.ReplaceAll(path1, `\`, "/")
		datFile, errData = os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
		if errData != nil {
			panic(errData)
		}
	}
	indFile, errInd := os.OpenFile(path+"/SSTableData/IndexFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if errInd != nil {
		indFile, errInd = os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
		if errInd != nil {
			panic(errInd)
		}
	}
	sumFile, errSum := os.OpenFile(path+"/SSTableData/SummaryFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if errSum != nil {
		sumFile, errSum = os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
		if errSum != nil {
			panic(errData)
		}
	}
	var offsetInd uint64 = 0
	stringovi := []string{}
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

		stringovi = append(stringovi, cvor.GetKey())

		size := 4 + 8 + 1 + 8 + 8 + key_u + val_u
		offset_ind := make([]byte, 8)
		binary.LittleEndian.PutUint64(offset_ind, offsetInd)

		indFile.Write(keySize)
		indFile.Write([]byte(cvor.GetKey()))
		indFile.Write(offset_ind)

		offsetInd = offsetInd + size

	}
	niz := MerkleTree.Pretvori_u_bajtove(stringovi)
	MerkleTree.Kreiraj_MerkleTree(niz, path+"/SSTableData/MetadataL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".txt")
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
func NadjiIndex(f *os.File, kljuc string, ima bool) (bool, uint64) {
	if ima {
		f.Seek(0, 0) // ukoliko je ima true, onda se pravi vise datoteka
	}
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

func NapraviSSTableJedanFajl(lCvor []*SkipList.SkipListNode, level int, index int) {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	var datFile *os.File
	datFile, errData := os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if errData != nil {
		path1, _ = filepath.Abs("../Projekat/Spojeno/Data")
		path = strings.ReplaceAll(path1, `\`, "/")
		datFile, errData = os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
		if errData != nil {
			panic(errData)
		}
	}
	var offsetInd int = 0
	stringovi := []string{}

	first := make([]byte, 8)
	first_u := uint64(len(lCvor[0].GetKey()))
	binary.LittleEndian.PutUint64(first, first_u)

	last := make([]byte, 8)
	last_u := uint64(len(lCvor[len(lCvor)-1].GetKey()))
	binary.LittleEndian.PutUint64(last, last_u)
	//for _, cvor := range lCvor {
	//	BloomFilter.AddNovi(cvor.GetKey(), path+"/SSTableData/AllDataFileL"+strconv.Itoa(level)+"Id"+strconv.Itoa(index)+".db")
	//}
	datFile.Close()
	datFile, errData = os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if errData != nil {
		panic(errData)
	}
	//bajt := make([]byte, 1)
	//velicinaBloom := 0
	//for {
	//	_, err := datFile.Read(bajt)
	//	if err != nil {
	//		break
	//	}
	//	velicinaBloom++
	//}
	//fmt.Println(velicinaBloom)
	//velicina, _ := datFile.Stat()
	//velicinaBloom := velicina.Size()
	//podaci := make([]byte, velicinaBloom)
	//datFile.Read(podaci)
	//datFile.Truncate(int64(velicinaBloom))
	//velB := make([]byte, 8)
	//binary.BigEndian.PutUint64(velB, uint64(velicinaBloom))
	//datFile.Write(velB)   //upisana velicina blooma
	//datFile.Write(podaci) //podaci iz blooma

	datFile.Write(first) //summary
	datFile.Write([]byte(lCvor[0].GetKey()))
	datFile.Write(last)
	datFile.Write([]byte(lCvor[len(lCvor)-1].GetKey()))

	velcinaInd := 0
	for _, cvor := range lCvor {
		velcinaInd += 8 + len(cvor.GetKey()) + 8
	}
	offsetInd += 8 + 8 + len(lCvor[0].GetKey()) + len(lCvor[len(lCvor)-1].GetKey()) + velcinaInd

	for _, cvor := range lCvor {

		keySize := make([]byte, 8)
		key_u := uint64(len(cvor.GetKey()))
		binary.LittleEndian.PutUint64(keySize, key_u)

		valSize := make([]byte, 8)
		val_u := uint64(len(cvor.GetValue()))
		binary.LittleEndian.PutUint64(valSize, val_u)

		size := 4 + 8 + 1 + 8 + 8 + key_u + val_u
		offset_ind := make([]byte, 8)
		binary.LittleEndian.PutUint64(offset_ind, uint64(offsetInd))

		datFile.Write(keySize)
		datFile.Write([]byte(cvor.GetKey()))
		datFile.Write(offset_ind)

		offsetInd = offsetInd + int(size)
	}
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
		stringovi = append(stringovi, cvor.GetKey())
	}

	niz := MerkleTree.Pretvori_u_bajtove(stringovi)
	MerkleTree.Kreiraj_MerkleTree(niz, path+"/SSTableData/MetadataL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index)+".txt")
	datFile.Close()

}
func NadjiBloom(kljuc string, f *os.File) bool {
	size := make([]byte, 8)
	f.Read(size)
	bloomSize := binary.LittleEndian.Uint64(size)
	bl := make([]byte, bloomSize)
	f.Read(bl) //procitani bajtovi blooma
	bloomFilter := BloomFilter.DeserijalizacijaNova(bl)
	nasao := bloomFilter.FindNovi(kljuc)
	return nasao

}
func NadjiSummaryJedanFajl(kljuc string, f *os.File) bool {
	size := make([]byte, 8)
	f.Read(size)
	keySize := binary.LittleEndian.Uint64(size) //dobijamo velicinu kljuca
	key_read := make([]byte, keySize)
	f.Read(key_read)
	key1 := string(key_read) //prvi kljuc
	fmt.Println(key1)
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

func NadjiIndexJedanFajl(f *os.File, kljuc string) (bool, uint64) {
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

func NadjiElementJedanFajl(offset uint64, f *os.File, kljuc string) (bool, []byte) {
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

func Svi_kljucevi_jednog_fajla(f *os.File) []string {
	f.Seek(0, 0)
	var kljucevi []string
	for true {
		size := make([]byte, 8)
		_, err := f.Read(size)
		if err != nil {
			return kljucevi
		}
		keySize := binary.LittleEndian.Uint64(size) //dobijamo velicinu kljuca
		key_read := make([]byte, keySize)
		f.Read(key_read)
		key := string(key_read)
		kljucevi = append(kljucevi, key)
		of := make([]byte, 8)
		f.Read(of)

	}
	return kljucevi
}

func Kompakcija(brojFajlova int, maxLevel int, level int, maxBloom int) {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	br := 0
	for br < brojFajlova-1 {
		skipList := SkipList.NapraviSkipList(10) //velicina koja
		br++
		bloom := BloomFilter.New_bloom(2*maxBloom, 0.05, level+1, br)
		l := strconv.Itoa(bloom.Level)
		i := strconv.Itoa(bloom.Index + 1)
		f1, err := os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(br)+".db", os.O_RDONLY, 0777)
		if err != nil {
			path1, _ = filepath.Abs("../Projekat/Spojeno/Data")
			path = strings.ReplaceAll(path1, `\`, "/")

			skipList = SkipList.NapraviSkipList(10) //velicina koja
			br++
			bloom := BloomFilter.New_bloom(2*maxBloom, 0.05, level+1, br)
			l = strconv.Itoa(bloom.Level)
			i = strconv.Itoa(bloom.Index + 1)
			f1, err = os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
				"Id"+strconv.Itoa(br)+".db", os.O_RDONLY, 0777)
			if err != nil {
				panic(err)
			}
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

		f1.Close()
		f2.Close()
		ObrisiFajlove(level, br)
		index := NovoImeFajla(level + 1)
		newData := skipList.GetElements()
		NapraviSSTable(newData, level+1, index)
		NapraviTOC(level+1, index)
	}
	index := NovoImeFajla(level + 1)
	if index == brojFajlova {
		if (level + 1) < maxLevel {
			Kompakcija(brojFajlova, maxLevel, level+1, 2*maxBloom)
		}
	}
}
func NapraviTOC(level int, index int) {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	tocFile, err := os.OpenFile(path+"/TOCFiles/TocFileL"+strconv.Itoa(level)+
		"Id"+strconv.Itoa(index-1)+".db", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		path1, _ = filepath.Abs("../Projekat/Spojeno/Data")
		path = strings.ReplaceAll(path1, `\`, "/")
		tocFile, err = os.OpenFile(path+"/TOCFiles/TocFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(index-1)+".db", os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			panic(err)
		}
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
		"Id" + strconv.Itoa(index-1) + ".txt"))
	if er != nil {
		panic(er)
	}
	_, er = tocFile.Write([]byte(path + "/SSTableData/MetaDataL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".txt"))
	if er != nil {
		panic(er)
	}
	tocFile.Close()

}

func ObrisiFajlove(level int, index int) {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	err1 := os.Remove(path + "/SSTableData/DataFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".db")
	if err1 != nil {
		path1, _ = filepath.Abs("../Projekat/Spojeno/Data")
		path = strings.ReplaceAll(path1, `\`, "/")
		err1 = os.Remove(path + "/SSTableData/DataFileL" + strconv.Itoa(level) +
			"Id" + strconv.Itoa(index) + ".db")
		if err1 != nil {
			panic(err1)
		}
	}
	fmt.Println(path + "/SSTableData/DataFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db")
	err2 := os.Remove(path + "/SSTableData/DataFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".db")
	if err2 != nil {
		panic(err1)
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
		"Id" + strconv.Itoa(index) + ".txt")
	if err7 != nil {
		panic(err7)
	}
	err8 := os.Remove(path + "/SSTableData/FilterFileL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index-1) + ".txt")
	if err8 != nil {
		panic(err8)
	}
	err9 := os.Remove(path + "/SSTableData/MetaDataL" + strconv.Itoa(level) +
		"Id" + strconv.Itoa(index) + ".txt")
	if err9 != nil {
		panic(err7)
	}
	err10 := os.Remove(path + "/SSTableData/MetaDataL" + strconv.Itoa(level) +
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
		f, err := os.OpenFile(path+"/SSTableData/DataFileL"+strconv.Itoa(level)+
			"Id"+strconv.Itoa(br)+".db", os.O_RDONLY, 0777)
		if err != nil { //nema takvog fajla, moze da bude novi
			f.Close()
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
