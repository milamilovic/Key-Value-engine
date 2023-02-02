package SimHash

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
)

// func main() {
// 	str1 := simHash("C:/Users/computer/Desktop/Faks/3. semestar/Napredni algoritmi i strukture podataka/Zadaci sa vezbi/vezbe3/simHash.txt")
// 	str2 := simHash("C:/Users/computer/Desktop/Faks/3. semestar/Napredni algoritmi i strukture podataka/Zadaci sa vezbi/vezbe3/simHash2.txt")
// 	fmt.Println(str1)
// 	fmt.Println(str2)
// 	fmt.Println(hamming(str1, str2))
// }

func SimHash(ime_fajla string) string {
	file, _ := os.OpenFile(ime_fajla, os.O_RDONLY, 0666)
	var reci []string
	Scanner := bufio.NewScanner(file)
	Scanner.Split(bufio.ScanWords)
	for Scanner.Scan() {
		reci = append(reci, Scanner.Text())
	}
	mapa := make(map[string]int, len(reci))
	for i := 0; i < len(reci); i++ {
		mapa[reci[i]] += 1
	}
	hesiranaMapa := make(map[string]string, len(mapa))
	for kljuc := range mapa {
		hesiranaMapa[kljuc] = hash(kljuc)
	}
	vrednosti := make([]int, len(hesiranaMapa[reci[0]]))
	for rec, hesirano := range hesiranaMapa {
		for i := range hesirano {
			mnozilac := int(hesirano[i] - 48)
			if mnozilac == 0 {
				mnozilac = -1
			}
			vrednosti[i] += mapa[rec] * mnozilac
		}
	}
	stringic := ""
	for i := range vrednosti {
		if vrednosti[i] <= 0 {
			stringic += "0"
		} else {
			stringic += "1"
		}
	}
	return stringic
}

func Hamming(str1 string, str2 string) int {
	num := 0
	for i := range str1 {
		if str1[i] != str2[i] {
			num += 1
		}
	}
	return num
}

func hash(kljuc string) string {
	return ToBinary(GetMD5Hash(kljuc))
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func ToBinary(s string) string {
	res := ""
	for _, c := range s {
		res = fmt.Sprintf("%s%.8b", res, c)
	}
	return res
}
