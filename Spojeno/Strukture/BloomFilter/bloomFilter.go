package BloomFilter

import (
	//"fmt"
	"encoding/gob"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// func main() {
// 	var kljucevi = [5]string{"cao", "mila", "nasp", "golang", "bloom"}
// 	var bloom = new_bloom(5, 0.01)
// 	var uspesno = true
// 	for i := 0; i < len(kljucevi); i++ {
// 		uspesno = bloom.add(kljucevi[i])
// 		if !uspesno {
// 			panic("doslo je do greske prilikom zapisivanja")
// 		}
// 	}
//  import "fmt"
// 	fmt.Println(bloom.find(kljucevi[3]))
// 	fmt.Println(bloom.find("lalala"))
// 	fmt.Println(bloom.find("kivi"))
// }

type BloomFilter struct {
	m                   uint
	k                   uint
	how_many_keys       int
	false_positive_rate float64
	Hashes              []Hash
	bytes               []byte
	Level               int
	Index               int
}

func New_bloom(how_many_keys int, false_positive float64, level int, index int) *BloomFilter {
	var m = CalculateM(how_many_keys, 0.01)
	var k = CalculateK(how_many_keys, m)
	hashes := Make_hashes(k, m)
	bytes := make([]byte, m)
	bloom := BloomFilter{m, k, how_many_keys, false_positive, hashes, bytes, level, index}
	Serijalizacija(&bloom)
	return &bloom
}

func Add(key string, filename string) bool {
	bloom := Deserijalizacija(filename)
	bytes := bloom.bytes

	for j := 0; j < len(bloom.Hashes); j++ {
		bytes[bloom.Hashes[j].Hash(key, int(bloom.m))] = 1
	}
	bloom.bytes = bytes
	Serijalizacija(bloom)
	return true
}
func AddNovi(key string, filename string) bool {
	bloom := Deserijalizacija(filename)
	bytes := bloom.bytes

	for j := 0; j < len(bloom.Hashes); j++ {
		bytes[bloom.Hashes[j].Hash(key, int(bloom.m))] = 1
	}
	bloom.bytes = bytes
	SerijalizacijaNova(bloom)
	return true
}

func Find(kljuc string, fajl string) bool {
	bloom := Deserijalizacija(fajl)
	bytes := bloom.bytes
	for j := 0; j < len(bloom.Hashes); j++ {
		if bytes[bloom.Hashes[j].Hash(kljuc, int(bloom.m))] == 0 {
			return false
		}
	}
	return true
}
func (bl *BloomFilter) FindNovi(kljuc string) bool {
	bytes := bl.bytes
	for j := 0; j < len(bl.Hashes); j++ {
		if bytes[bl.Hashes[j].Hash(kljuc, int(bl.m))] == 0 {
			return false
		}
	}
	return true
}

type Hash struct {
	Broj int
}

func Make_hashes(k uint, m uint) []Hash {
	hashes := make([]Hash, k)
	for i := 0; i < int(k); i++ {
		hashes[i] = Hash{rand.Intn(30492570)}
	}
	return hashes
}

func (h *Hash) Hash(kljuc string, m int) int {
	var vrednost = 0
	chars := []rune(kljuc)
	for i := 0; i < len(chars); i++ {
		vrednost += int(chars[i])
	}
	var hesirana = h.Broj * vrednost
	return hesirana % m
}

func Serijalizacija(bloom *BloomFilter) {

	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	l := strconv.Itoa(bloom.Level)
	i := strconv.Itoa(bloom.Index)
	file, err := os.OpenFile(path+"/SSTableData/filterFileL"+l+"Id"+i+".txt", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(&bloom)
	if err != nil {
		panic(err)
	}
}
func SerijalizacijaNova(bloom *BloomFilter) {

	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	l := strconv.Itoa(bloom.Level + 1)
	i := strconv.Itoa(bloom.Index + 1)
	file, err := os.OpenFile(path+"/SSTableData/AllDataFileL"+l+"Id"+i+".db", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(&bloom)
	if err != nil {
		panic(err)
	}
}

func Deserijalizacija(str string) *BloomFilter {
	// path1, _ := filepath.Abs("../Spojeno/Data")
	// path := strings.ReplaceAll(path1, `\`, "/")
	bloom := BloomFilter{}
	file, err := os.OpenFile(str, os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	decoder := gob.NewDecoder(file)
	_ = decoder.Decode(&bloom)
	// if err != nil {
	// 	panic(err)
	// }
	hashes := Make_hashes(bloom.k, bloom.m)
	bloom.Hashes = hashes
	return &bloom
}
func DeserijalizacijaNova(niz []byte) *BloomFilter {
	path1, _ := filepath.Abs("../Spojeno/Data")
	path := strings.ReplaceAll(path1, `\`, "/")
	bloom := BloomFilter{}
	f, err := os.OpenFile(path+"/SSTableData/BloomDes.db", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	f.Write(niz)
	decoder := gob.NewDecoder(f)
	_ = decoder.Decode(&bloom)
	// if err != nil {
	// 	panic(err)
	// }
	hashes := Make_hashes(bloom.k, bloom.m)
	bloom.Hashes = hashes
	return &bloom

}

func CalculateM(expectedElements int, falsePositiveRate float64) uint {
	return uint(math.Ceil(float64(expectedElements) * math.Abs(math.Log(falsePositiveRate)) / math.Pow(math.Log(2), float64(2))))
}

func CalculateK(expectedElements int, m uint) uint {
	return uint(math.Ceil((float64(m) / float64(expectedElements)) * math.Log(2)))
}
