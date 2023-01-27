package strukture

func main() {
	memTable := createMemTable(20, 15)
	memTable.add("1", []byte("a"))
	memTable.update("1", []byte("b"))
	memTable.deleteElement("1")
	flush := memTable.checkFlush()
	if flush == true {
		memTable.flush()
	}
}

type memTable struct {
	elementi         *SkipList
	velicina         int //velicina skipListe
	maxVelicina      int //maksimalna velicina za memTable
	trenutnaVelicina int
}

func createMemTable(max, velicina int) *memTable {
	elementi := makeSkipList(velicina)
	return &memTable{elementi, velicina, max, 0}
}

func (memTable *memTable) add(key string, value []byte) {
	b, cvor := memTable.elementi.findElement(key)
	if b == false {
		if cvor == nil {
			memTable.elementi.add(key, value)
			print("Ubacili smo novi element u skip listu")
			memTable.trenutnaVelicina++
		}
	}
}
func (memTable *memTable) update(key string, value []byte) {
	b, cvor := memTable.elementi.findElement(key)
	if b == true { //nasao je elemnt i menja mu value
		memTable.elementi.add(key, value)
		print("Izmenili smo element u skip listi")
	} else {
		if cvor != nil { //cvor je logicki obrisan
			memTable.elementi.add(key, value) //izmenice mu i tombstone na false
		}
	}
}
func (memTable *memTable) deleteElement(key string) {
	b, cvor := memTable.elementi.findElement(key)
	if b == true { //nasao je elemnt i menja mu value
		memTable.elementi.logDelete(key)
		print("Izbrisali smo element u skip listi")
	} else {
		if cvor != nil { //cvor je logicki obrisan
			print("Element je vec logicki obrisan")
		} else {
			print("Nema element sa unetim kljucem")
		}
	}
}

func (memTable *memTable) checkFlush() bool {
	if memTable.maxVelicina <= memTable.trenutnaVelicina {
		return true //treba flush odraditi
	} else {
		return false
	}
}

func (memTable *memTable) flush() {
	memTable.writeSSTable()
	memTable = createMemTable(15, 20) //pre ovoga treba upisati na disk, SStable
}

func (memTable *memTable) writeSSTable() {
	return
}
