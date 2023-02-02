package dodavanje

import (
	"Strukture/MemTableBTree"
	"Strukture/MemTableSkipList"
	"Strukture/Wal"
)

func dodaj_skiplist(key string, value []byte, mt *MemTableSkipList.MemTable, w *Wal.Wal) {
	w.Dodaj_u_wal(key, value, false) // tombstone je false kada dodajemo
	mt.Add(key, value)
	if mt.ProveriFlush() {
		mt.Flush()
	}
}

func dodaj_bstablo(key string, value []byte, mt *MemTableBTree.MemTable, w *Wal.Wal) {
	w.Dodaj_u_wal(key, value, false) // tombstone je false kada dodajemo
	mt.Add(key, value)
	if mt.ProveriFlush() {
		mt.Flush()
	}
}
