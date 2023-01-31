package main

import (
	"time"
)

// func main() {
// 	tb := newTokenBucket(9999999, 2, 5)
// 	fmt.Println(tb.check(9999999))
// 	fmt.Println(tb.check(9999999))
// 	fmt.Println(tb.check(9999999))
// 	fmt.Println(tb.check(9999999))
// 	vreme := time.Now().Unix()
// 	fmt.Println("cekamo...")
// 	for time.Now().Unix()-vreme < 6 {
// 	}
// 	fmt.Println(tb.check(9999999))
// 	fmt.Println(tb.check(9999999))
// 	fmt.Println(tb.check(9999999))
// }

type TokenBucket struct {
	kljuc               int //podrazumevamo da imamo jednog korisnika
	broj_tokena         int
	max_tokena          int
	poslednji_timestamp int64
	interval            int64 //sekunde
}

func newTokenBucket(key int, maxtok int, interv int64) *TokenBucket {
	return &TokenBucket{key, maxtok, maxtok, time.Now().Unix(), interv}
}

func (tokenBucket *TokenBucket) check(kljuc int) bool {
	if kljuc != tokenBucket.kljuc {
		return false
	}
	//ako je interval istekao resetujemo
	if time.Now().Unix()-tokenBucket.poslednji_timestamp > tokenBucket.interval {
		tokenBucket.poslednji_timestamp = time.Now().Unix()
		tokenBucket.broj_tokena = tokenBucket.max_tokena
	}
	//proveravamo da li imamo dovoljno tokena
	if tokenBucket.broj_tokena > 0 {
		tokenBucket.broj_tokena--
		return true
	}
	return false
}
