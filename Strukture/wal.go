package main

import (
	"encoding/binary"
	"hash/crc32"
	"os"
	"strconv"
)

// func main() {
// 	wal := napraviWal("", 0)
// 	wal.dodaj_u_wal("5", []byte("vrednost od 5"), false)
// 	wal.dodaj_u_wal("7", []byte("a"), false)
// 	wal.dodaj_u_wal("8", []byte("a"), false)
// 	wal.dodaj_u_wal("334", []byte("a"), false)
// 	wal.dodaj_u_wal("sdkfj", []byte("a"), false)
// 	wal.dodaj_u_wal("fff", []byte("vrednost od fff"), false)
// 	wal.dodaj_u_wal("12", []byte("vrednost od 12"), false)
// 	wal.dodaj_u_wal("9", []byte("vrednost od 9"), false)
// 	fmt.Println("Trazimo 9")
// 	bajtovi := wal.nadji_podatak("9")
// 	fmt.Println("vrednost u bajtovima: ")
// 	fmt.Println(bajtovi)
// 	fmt.Println("dekodirana vrednost: ")
// 	fmt.Println(string(bajtovi))
// 	fmt.Println("Trazimo fff")
// 	bajtovi = wal.nadji_podatak("fff")
// 	fmt.Println("vrednost u bajtovima: ")
// 	fmt.Println(bajtovi)
// 	fmt.Println("dekodirana vrednost: ")
// 	fmt.Println(string(bajtovi))
// 	fmt.Println("Trazimo 12")
// 	bajtovi = wal.nadji_podatak("12")
// 	fmt.Println("vrednost u bajtovima: ")
// 	fmt.Println(bajtovi)
// 	fmt.Println("dekodirana vrednost: ")
// 	fmt.Println(string(bajtovi))
// }

/*
   +---------------+-----------------+---------------+---------------+-----------------+-...-+--...--+
   |    CRC (4B)   | Timestamp (8B) | Tombstone(1B) | Key Size (8B) | Value Size (8B) | Key | Value |
   +---------------+-----------------+---------------+---------------+-----------------+-...-+--...--+
   CRC = 32bit hash computed over the payload using CRC
   Key Size = Length of the Key data
   Tombstone = If this record was deleted and has a value
   Value Size = Length of the Value data
   Key = Key data
   Value = Value data
   Timestamp = Timestamp of the operation in seconds
*/

const (
	CRC_SIZE        = 4
	TIMESTAMP_SIZE  = 8
	TOMBSTONE_SIZE  = 1
	KEY_SIZE_SIZE   = 8
	VALUE_SIZE_SIZE = 8

	CRC_START        = 0
	TIMESTAMP_START  = CRC_START + CRC_SIZE
	TOMBSTONE_START  = TIMESTAMP_START + TIMESTAMP_SIZE
	KEY_SIZE_START   = TOMBSTONE_START + TOMBSTONE_SIZE
	VALUE_SIZE_START = KEY_SIZE_START + KEY_SIZE_SIZE
	KEY_START        = VALUE_SIZE_START + VALUE_SIZE_SIZE
)

func CRC32(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

type Segment struct {
	index        int
	podaci       []byte
	velicina     int
	max_velicina int
}

func (segment *Segment) nabavi_podatke() []byte {
	return segment.podaci
}

func (segment *Segment) add(podatak []byte) int {
	//upisujemo podatke dok mozemo
	//ako ne moze sve da stane vracamo false, a ako moze true
	if segment.velicina+len(podatak) > segment.max_velicina {
		//vracamo koliko moze da stane
		return segment.max_velicina - segment.velicina
	}
	segment.velicina += len(podatak)
	for i := 0; i < len(podatak); i++ {
		segment.podaci = append(segment.podaci, podatak[i])
	}
	return -1
}

func (segment *Segment) zapisi(putanja string) {
	//kreramo novi fajl sa zadatom putanjom
	naziv_fajla := putanja + "wal_" + strconv.FormatInt(int64(segment.index), 10) + ".bin"
	file, err := os.OpenFile(naziv_fajla, os.O_CREATE, 0666)
	//zapisemo podatke
	err = binary.Write(file, binary.LittleEndian, segment.podaci)
	if err != nil {
		panic(err)
	}
	//ne smemo da zaboravimo da zatvorimo fajl!!!
	file.Close()
}

type Wal struct {
	putanja                       string
	velicina_segmenta             int
	redni_broj_trenutnog_segmenta int
	low_water_mark                int
	segmenti                      []*Segment
	trenutni_segment              *Segment
	imena_segmenata_po_indeksu    map[int]string
	podaci_iz_wal                 map[string][]byte
}

func napraviWal(putanja string, low_water_mark int) *Wal {
	wal := Wal{putanja, 50, 0, low_water_mark, make([]*Segment, 0), &Segment{0, make([]byte, 0, 50), 0, 50}, make(map[int]string), make(map[string][]byte)}
	wal.segmenti = append(wal.segmenti, wal.trenutni_segment)
	return &wal
}

func (wal *Wal) novi_segment() {
	//pravimo novi segment, dodajemo ga u listu segmenata
	novi := Segment{wal.redni_broj_trenutnog_segmenta + 1, make([]byte, 0, 50), 0, 50}
	wal.segmenti = append(wal.segmenti, &novi)
	//zapisujemo stari
	wal.trenutni_segment.zapisi(wal.putanja)
	//a novi postavljamo kao trenutni
	wal.trenutni_segment = &novi
	wal.redni_broj_trenutnog_segmenta += 1
	wal.imena_segmenata_po_indeksu[wal.redni_broj_trenutnog_segmenta] = wal.putanja + "wal_" + strconv.FormatInt(int64(novi.index), 10) + ".bin"
}

func (wal *Wal) dodaj_u_wal(key string, value []byte, tombstone_bool bool) {
	//prvo treba da napravimo da podaci budu u odgovarajucem formatu
	//CRC (4B)   | Timestamp (8B) | Tombstone (1B) | Key Size (8B) | Value Size (8B) | Key | Value
	crc := make([]byte, 4)
	binary.LittleEndian.PutUint32(crc, crc32.ChecksumIEEE(value))
	timestamp := make([]byte, 8)
	tombstone := make([]byte, 1)
	if tombstone_bool {
		tombstone[0] = 1
	} else {
		tombstone[0] = 0
	}
	key_size := make([]byte, 8)
	binary.LittleEndian.PutUint64(key_size, uint64(len(key)))
	value_size := make([]byte, 8)
	binary.LittleEndian.PutUint64(value_size, uint64(len(value)))
	//pa onda da objedinimo te podatke
	podaci := make([]byte, 0)
	podaci = append(podaci, crc...)
	podaci = append(podaci, timestamp...)
	podaci = append(podaci, tombstone...)
	podaci = append(podaci, key_size...)
	podaci = append(podaci, value_size...)
	podaci = append(podaci, []byte(key)...)
	podaci = append(podaci, []byte(value)...)
	//i dodamo ih u segment
	koliko_je_dodato := wal.trenutni_segment.add(podaci)
	pomeraj := 0
	//ako popunimo segment treba da napravimo novi i u njega zapisemo ostale podatke
	for koliko_je_dodato > 0 {
		pomeraj += koliko_je_dodato
		wal.novi_segment()
		koliko_je_dodato = wal.trenutni_segment.add(podaci[pomeraj:])
	}
	wal.podaci_iz_wal[key] = value
}

func (wal *Wal) nadji_podatak(kljuc string) []byte {
	podatak := wal.podaci_iz_wal[kljuc]
	return podatak
}

func (wal *Wal) brisi_segmente() {
	//obrisemo segmente ispod watermarka
	for i := 0; i < wal.low_water_mark; i++ {
		err := os.Remove(wal.putanja + "wal_" + strconv.FormatInt(int64(i), 10) + ".bin")
		if err != nil {
			panic(err)
		}
		delete(wal.imena_segmenata_po_indeksu, i)
	}
	//i onda preimenujemo ostale
	stari_fajlovi, _ := os.ReadDir(wal.putanja)
	for id, fajl_podaci := range stari_fajlovi {
		staro_ime := wal.putanja + fajl_podaci.Name()
		novo_ime := wal.putanja + "wal_" + strconv.FormatInt(int64(id+1), 10) + ".bin"
		err := os.Rename(staro_ime, novo_ime)
		if err != nil {
			panic(err)
		}
	}
}

func (wal *Wal) sacuvaj() {
	wal.trenutni_segment.zapisi(wal.putanja)
}
