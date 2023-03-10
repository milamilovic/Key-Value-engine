package BloomFilter

import (
	//"fmt"
	"math"
	"math/rand"
	"os"
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
}

func New_bloom(how_many_keys int, false_positive float64) *BloomFilter {
	var m = CalculateM(how_many_keys, 0.01)
	var k = CalculateK(how_many_keys, m)
	hashes := Make_hashes(k, m)
	file, _ := os.OpenFile("C:/Users/Sonja/Desktop/Key-Value-engine/Strukture/BloomFilter/filter.txt", os.O_CREATE, 0666)
	bytes := make([]byte, m)
	_, _ = file.WriteAt(bytes, 0)
	file.Close()
	return &BloomFilter{m, k, how_many_keys, false_positive, hashes, bytes}
}

func (bloom *BloomFilter) Add(key string) bool {
	file, _ := os.OpenFile("C:/Users/Sonja/Desktop/Key-Value-engine/Strukture/BloomFilter/filter.txt", os.O_CREATE, 0666)
	bytes := make([]byte, int(bloom.m))
	_, err := file.Read(bytes)
	if err != nil {
		panic(err)
	}
	for j := 0; j < len(bloom.Hashes); j++ {
		bytes[bloom.Hashes[j].Hash(key, int(bloom.m))] = 1
	}
	_, err = file.WriteAt(bytes, 0)
	file.Close()
	if err != nil {
		return false
	} else {
		return true
	}
}

func (bloom *BloomFilter) Find(kljuc string) bool {
	file, _ := os.OpenFile("C:/Users/Sonja/Desktop/Key-Value-engine/Strukture/BloomFilter/filter.txt", os.O_CREATE, 0666)
	bytes := bloom.bytes
	_, _ = file.Read(bytes)
	for j := 0; j < len(bloom.Hashes); j++ {
		if bytes[bloom.Hashes[j].Hash(kljuc, int(bloom.m))] == 0 {
			return false
		}
	}
	return true
}

type Hash struct {
	broj int
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
	var hesirana = h.broj * vrednost
	return hesirana % m
}

func CalculateM(expectedElements int, falsePositiveRate float64) uint {
	return uint(math.Ceil(float64(expectedElements) * math.Abs(math.Log(falsePositiveRate)) / math.Pow(math.Log(2), float64(2))))
}

func CalculateK(expectedElements int, m uint) uint {
	return uint(math.Ceil((float64(m) / float64(expectedElements)) * math.Log(2)))
}
