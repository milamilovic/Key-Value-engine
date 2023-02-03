package dodavanje

import (
	"Strukture/BloomFilter"
	"Strukture/MemTableBTree"
	"Strukture/MemTableSkipList"
	"Strukture/TokenBucket"
	"Strukture/Wal"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func Dodaj_skiplist(key string, value []byte, mt *MemTableSkipList.MemTable, w *Wal.Wal, t *TokenBucket.TokenBucket, bloom *BloomFilter.BloomFilter) {
	if t.Check(t.Kljuc) {
		w.Dodaj_u_wal(key, value, false) // tombstone je false kada dodajemo
		mt.Add(key, value)
		l := strconv.Itoa(bloom.Level)
		i := strconv.Itoa(bloom.Index)
		path1, _ := filepath.Abs("../Spojeno/Data")
		path := strings.ReplaceAll(path1, `\`, "/")
		BloomFilter.Add(key, path+"/SSTableData/filterFileL"+l+"Id"+i+".txt")
	} else {
		fmt.Println("Neuspesno dodavanje, isteklo je vreme.")
	}
	//w.Dodaj_u_wal(key, value, false) // tombstone je false kada dodajemo
	//mt.Add(key, value)
	// if mt.ProveriFlush() {
	// 	mt.Flush()
	// }
}

func Dodaj_bstablo(key string, value []byte, mt *MemTableBTree.MemTable, w *Wal.Wal, t *TokenBucket.TokenBucket) {
	if t.Check(t.Kljuc) {
		w.Dodaj_u_wal(key, value, false) // tombstone je false kada dodajemo
		mt.Add(key, value)
	} else {
		fmt.Println("Neuspesno dodavanje, isteklo je vreme.")
	}
	// if mt.ProveriFlush() {
	// 	mt.Flush()
	// }
}
