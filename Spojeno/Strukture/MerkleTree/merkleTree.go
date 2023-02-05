package MerkleTree

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
)

type Hash [20]byte //vrednost

func (n Node) String() string {
	return hex.EncodeToString(n.data[:])
}

func hash(podaci []byte) Hash {
	return sha1.Sum(podaci)
}

type Root struct {
	root *Node
}

func (mr *Root) String() string {
	return mr.root.String()
}

type Node struct {
	data  [20]byte
	levi  *Node
	desni *Node
}

type PrazanNode struct {
	//data  [20]byte
	//levi  nil
	//desni nil
}

func Pretvori_u_bajtove(stringovi []string) [][]byte {
	data := [][]byte{}
	for i := 0; i < len(stringovi); i++ {
		key_byte := []byte(stringovi[i])
		//fmt.Println(key_byte)
		data = append(data, key_byte)
	}
	return data
}

func Kreiraj_listove(data [][]byte) []*Node {
	listovi := []*Node{}

	for i := 0; i < len(data); i++ {
		node := Node{hash(data[i]), nil, nil}
		listovi = append(listovi, &node)
	}

	return listovi
}

// funkcija za kreiranje merkl stabla, "pocetna" funkcija
func Kreiraj_MerkleTree(keys [][]byte, putanja string) *Root {

	data := keys

	listovi := Kreiraj_listove(data)
	root_node := Kreiraj_cvorove(listovi)

	root := Root{root_node}
	Upisi_u_fajl(root_node, putanja)
	return &root
}

func Kreiraj_cvorove(svi_listovi []*Node) *Node {
	nivo := []*Node{}
	cvorovi := svi_listovi
	if len(cvorovi) == 1 {
		// samo koren
		return cvorovi[0]
	}
	if len(cvorovi) > 1 {
		for i := 0; i < len(cvorovi); i += 2 {
			if (i + 1) < len(cvorovi) {
				prvi := cvorovi[i]
				drugi := cvorovi[i+1]
				novi_cvor_podaci := append(prvi.data[:], drugi.data[:]...)
				novi_cvor := Node{hash(novi_cvor_podaci), prvi, drugi}
				nivo = append(nivo, &novi_cvor)
			} else {
				prvi := cvorovi[i]
				drugi := Node{data: [20]byte{}, levi: nil, desni: nil} // ako je jedan cvor visak, njega formiramo kao prazan cvor
				novi_cvor_podaci := append(prvi.data[:], drugi.data[:]...)
				novi_cvor := Node{hash(novi_cvor_podaci), prvi, &drugi}
				nivo = append(nivo, &novi_cvor)
			}
			//nivo := append(nivo, &novi_cvor)
		}
		cvorovi = nivo

		if len(cvorovi) == 1 {
			return cvorovi[0]
		}
	}
	return Kreiraj_cvorove(nivo)
}

func PrintTree(root *Node) {
	queue := make([]*Node, 0)
	queue = append(queue, root)

	for len(queue) != 0 {
		e := queue[0]
		queue = queue[1:]
		//fmt.Println(e.String())

		if e.levi != nil {
			queue = append(queue, e.levi)
		}
		if e.desni != nil {
			queue = append(queue, e.desni)
		}
	}
}

func Upisi_u_fajl(root *Node, naziv_fajla string) {
	file, err := os.Create(naziv_fajla)

	if err != nil {
		log.Fatalf("Neuspesno kreiranje fajla: %s", err)
	}
	defer file.Close()
	queue := make([]*Node, 0)
	queue = append(queue, root)

	for len(queue) != 0 {
		e := queue[0]
		queue = queue[1:]
		_, err = file.WriteString(e.String() + "\n")
		if e.levi != nil {
			queue = append(queue, e.levi)
		}
		if e.desni != nil {
			queue = append(queue, e.desni)
		}
	}
	if err != nil {
		log.Fatalf("Neuspesno upisivanje u fajl: %s", err)
	}
}
