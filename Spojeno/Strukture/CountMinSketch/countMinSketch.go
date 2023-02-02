package CountMinSketch

import (
	"math"
	"math/rand"
	"os"
)

// func main() {
// 	CountMin := createCMS(0.1, 0.1)
// 	file, _ := os.OpenFile("C:/Users/Sonja/Desktop/napredni.txt", os.O_CREATE, 0666)
// 	for i := 0; i < int(CountMin.k); i++ {
// 		z := int64(i * int(CountMin.m))
// 		_, _ = file.WriteAt(CountMin.bytes, z)
// 	}
// 	file.Close()
// 	var kljucevi = []string{"cao", "mila", "nasp", "golang", "bloom", "cao", "cao", "nasp", "golang"}
// 	var uspesno = true
// 	for i := 0; i < len(kljucevi); i++ {
// 		uspesno = CountMin.add(kljucevi[i], CountMin.hashes, int(CountMin.m))
// 		if !uspesno {
// 			panic("doslo je do greske prilikom zapisivanja")
// 		}
// 	}

// 	fmt.Println(CountMin.cms(kljucevi[0], CountMin.hashes, int(CountMin.m)))
// 	fmt.Println(CountMin.cms(kljucevi[3], CountMin.hashes, int(CountMin.m)))
// 	fmt.Println(CountMin.cms("kivi", CountMin.hashes, int(CountMin.m)))
// }

type CountMinSketch struct {
	k       uint
	M       uint
	Hashes  []Hash
	delta   float64
	epsilon float64
	Bytes   []byte
}

func CreateCMS(epsilon, delta float64) *CountMinSketch {
	var m = CalculateM(epsilon)
	var k = CalculateK(delta)
	hashes := HashFunctions(k)
	bytes := make([]byte, m)
	return &CountMinSketch{k: k, M: m, Hashes: hashes, delta: delta, epsilon: epsilon, Bytes: bytes}

}

func HashFunctions(k uint) []Hash {
	hashes := make([]Hash, k)
	for i := 0; i < int(k); i++ {
		hashes[i] = Hash{rand.Intn(30492570)}
	}
	return hashes
}

func CalculateM(epsilon float64) uint {
	return uint(math.Ceil(math.E / epsilon))
}

func CalculateK(delta float64) uint {
	return uint(math.Ceil(math.Log(math.E / delta)))
}

type Hash struct {
	broj int
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

func Add(key string, hashes []Hash, m int, bajtovi []byte) []byte {
	for i := 0; i < int(len(hashes)); i++ {
		bajtovi[hashes[i].Hash(key, m)] += 1
	}
	return bajtovi
}

func Cms(kljuc string, hashes []Hash, m int, bajtovi []byte) int {
	min := 100
	for i := 0; i < int(len(hashes)); i++ {
		for j := 0; j < len(hashes); j++ {
			if bajtovi[hashes[j].Hash(kljuc, m)] < byte(min) {
				min = int(bajtovi[hashes[j].Hash(kljuc, m)])
			}
		}
	}
	return min
}

func Serijalizacija(bajtovi []byte, filename string) {
	os.WriteFile(filename, bajtovi, 0666)
}
