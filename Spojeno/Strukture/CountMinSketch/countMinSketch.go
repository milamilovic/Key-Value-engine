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
	m       uint
	hashes  []Hash
	delta   float64
	epsilon float64
	bytes   []byte
}

func CreateCMS(epsilon, delta float64) *CountMinSketch {
	var m = CalculateM(epsilon)
	var k = CalculateK(delta)
	hashes := make([]Hash, k)
	hashes = HashFunctions(k)
	file, _ := os.OpenFile("C:/Users/Sonja/Desktop/napredni.txt", os.O_CREATE, 0666)
	bytes := make([]byte, m)
	_, _ = file.WriteAt(bytes, 0)
	file.Close()

	return &CountMinSketch{k: k, m: m, hashes: hashes, delta: delta, epsilon: epsilon, bytes: bytes}

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

func (countMin *CountMinSketch) Add(key string, hashes []Hash, m int) bool {
	file, err := os.OpenFile("C:/Users/Sonja/Desktop/napredni.txt", os.O_CREATE, 0666)
	for i := 0; i < int(len(hashes)); i++ {
		z := int64(i * int(m))
		bytes := make([]byte, m)
		file.ReadAt(bytes, z)
		for j := 0; j < len(hashes); j++ {
			bytes[hashes[j].Hash(key, m)] += 1
		}
		_, err = file.WriteAt(bytes, z)
	}
	file.Close()
	if err != nil {
		return false
	} else {
		return true
	}
}

func (countMin *CountMinSketch) Cms(kljuc string, hashes []Hash, m int) int {
	file, _ := os.OpenFile("C:/Users/Sonja/Desktop/napredni.txt", os.O_CREATE, 0666)
	min := 100
	for i := 0; i < int(len(hashes)); i++ {
		z := int64(i * int(m))
		bytes := make([]byte, m)
		file.ReadAt(bytes, z)
		for j := 0; j < len(hashes); j++ {
			if bytes[hashes[j].Hash(kljuc, m)] < byte(min) {
				min = int(bytes[hashes[j].Hash(kljuc, m)])
			}
		}
	}
	return min
}
