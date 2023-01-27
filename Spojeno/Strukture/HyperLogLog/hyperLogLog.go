package HyperLogLog

import (
	"bytes"
	"encoding/gob"
	"math"
	"strconv"
)

// func main() {
// 	var kljucevi = [12]string{"cao", "mila", "nasp", "golang", "hll", "kivi", "kruska", "probabilisticka struktura", "kreativnost", "sistematicno", "hihi", "vau"}
// 	var hyper = makeHyperLogLog(HLL_MAX_PRECISION)
// 	var uspesno = true
// 	for i := 0; i < len(kljucevi); i++ {
// 		uspesno = hyper.add(kljucevi[i])
// 		if !uspesno {
// 			panic("doslo je do greske prilikom zapisivanja")
// 		}
// 	}
// 	fmt.Println(hyper.Estimate())
// }

const (
	HLL_MIN_PRECISION = 4
	HLL_MAX_PRECISION = 16
)

func makeHyperLogLog(preciznost int) HLL {
	var m = uint64(math.Pow(2, float64(preciznost)))
	return HLL{m: m, p: uint8(preciznost), reg: make([]uint8, m)}
}

type HLL struct {
	m   uint64
	p   uint8
	reg []uint8
}

func (hyper *HLL) add(key string) bool {
	//hesiran kljuc je string od binarnog broja
	hesiran_kljuc := hyper.hesiraj(key)
	baket, err := strconv.ParseInt(hesiran_kljuc[0:hyper.p], 2, 0)
	broj_vodecih_nula := 0
	for i := 0; i < len(hesiran_kljuc)-int(hyper.p); i++ {
		if hesiran_kljuc[len(hesiran_kljuc)-1-i] == 48 {
			broj_vodecih_nula++
		} else {
			break
		}
	}
	//fmt.Println(broj_vodecih_nula)
	if hyper.reg[baket] < uint8(broj_vodecih_nula) {
		hyper.reg[baket] = uint8(broj_vodecih_nula)
	}
	if err == nil {
		return true
	}
	return false
}

func (hyper *HLL) hesiraj(kljuc string) string {
	var vrednost = 0
	chars := []rune(kljuc)
	for i := 0; i < len(chars); i++ {
		vrednost += int(chars[i])
	}
	var hesirana = (vrednost*12348912734 + 934738) % 67280421310721
	string_resenja := strconv.FormatInt(int64(hesirana), 2)
	if len(string_resenja) > 32 {
		string_resenja = string_resenja[len(string_resenja)-32:]
	}
	if len(string_resenja) < 32 {
		hesirana *= 67280421310721
		string_resenja := strconv.FormatInt(int64(hesirana), 2)
		if len(string_resenja) > 32 {
			string_resenja = string_resenja[len(string_resenja)-32:]
		}
	}
	return string_resenja
}

// procenjuje koliko ima elemenata u hll-u
func (hll *HLL) Estimate() float64 {
	sum := 0.0
	for _, val := range hll.reg {
		sum += math.Pow(math.Pow(2.0, float64(val)), -1)
	}

	alpha := 0.7213 / (1.0 + 1.079/float64(hll.m))
	estimation := alpha * math.Pow(float64(hll.m), 2.0) / sum
	emptyRegs := hll.emptyCount()
	if estimation <= 2.5*float64(hll.m) { // do small range correction
		if emptyRegs > 0 {
			estimation = float64(hll.m) * math.Log(float64(hll.m)/float64(emptyRegs))
		}
	} else if estimation > 1/30.0*math.Pow(2.0, 32.0) { // do large range correction
		estimation = -math.Pow(2.0, 32.0) * math.Log(1.0-estimation/math.Pow(2.0, 32.0))
	}
	return estimation
}

// procenjuje koliko ima praznih baketa
func (hll *HLL) emptyCount() int {
	sum := 0
	for _, val := range hll.reg {
		if val == 0 {
			sum++
		}
	}
	return sum
}

func deserijalizacija(podaci []byte) *HLL {
	bajtovi := bytes.NewBuffer(podaci)
	dekoder := gob.NewDecoder(bajtovi)
	hyper := new(HLL)
	for {
		err := dekoder.Decode(&hyper)
		if err != nil {
			break
		}
	}
	return hyper
}

func serijalizacija(hyper *HLL) []byte {
	var podaci bytes.Buffer
	koder := gob.NewEncoder(&podaci)
	koder.Encode(&hyper)
	return podaci.Bytes()
}
