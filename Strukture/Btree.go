package main

func main() {
	btree := makeBtree(3)
	btree.add("3", []byte("a"))
	btree.add("5", []byte("a"))
	btree.add("9", []byte("a"))
	btree.add("1", []byte("a"))
	btree.add("7", []byte("a"))
	btree.add("1", []byte("ponovo isti kljuc"))
	btree.add("4", []byte("a"))
	b, _ := btree.findElement("1")
	print("\n", b)
	b, _ = btree.findElement("5")
	print("\n", b)
	b, _ = btree.findElement("7")
	print("\n", b)
	b, _ = btree.findElement("12") //ovo je false
	print("\n", b)

	btree.delete("4")
	b, _ = btree.findElement("4") //false jer je logicki obrisan
	print("\n", b, "\n")
	btree.delete("4") //brisemo ga opet

	btree.add("4", []byte("a"))   //postoji ali je obrisan pa ga azuriramo
	b, _ = btree.findElement("4") //azurirano je pa je true
	print("\n", b)

}

type Btree struct {
	root      BtreeNode
	maxHeight int
	height    int
}

type BtreeNode struct {
	elements     []*Node
	children     []*BtreeNode
	max_children int
}

type Node struct {
	key   string
	value []byte
}

func makeBtree(maxHeight int) *Btree {
	root := BtreeNode{make([]*Node, 3), make([]*BtreeNode, 3), 1}
	return &Btree{root, maxHeight + 1, 1}
}

func (btree *Btree) findElement(key string) (bool, *Node) {
	trenutniCvor := btree.root
	k := 0
	for i := 0; i <= btree.height; i++ {
		trenutnoDete := trenutniCvor.elements[k]
		sledeceDete := trenutniCvor.elements[k+1]
		for sledeceDete != nil {
			//ako je ovaj dete tjt nasli smo
			if trenutnoDete.key == key {
				return true, trenutnoDete
			}
			//ako je sledeci dete tjt nasli smo
			if sledeceDete.key == key {
				return true, sledeceDete
			}
			//ako je kljuc manji od trenutnog kljuca idemo na dete tog indeksa
			if key < trenutnoDete.key {
				trenutniCvor = *trenutniCvor.children[k]
				k = 0
				continue
			}
			//ako je kljuc izmedju trenutnog i sledeceg idemo tamo
			if sledeceDete.key > key && trenutnoDete.key < key {
				trenutniCvor = *trenutniCvor.children[k+1]
				k = 0
				continue
			}
			k++
		}
		if trenutniCvor.max_children == 0 {
			return false, sledeceDete
		}
		//ako je sl dete nil onda ili nema deteta ili je ,,desno"
		if key > trenutnoDete.key {
			trenutniCvor = *trenutniCvor.children[k+1]
			k = 0
		}
	}
	return false, nil //nismo upste nasli
}

func (btree *Btree) add(key string, value []byte) {
	b, _ := btree.findElement(key)
	trenutniCvor := btree.root
	roditelj := trenutniCvor
	if b == false { //ovo znaci da ne postoji jer ga nije pronasao
		pronadjen := false
		//idemo sve do lista
		for trenutniCvor.children[0] != nil {
			for i := 0; i < len(trenutniCvor.elements)-1 && trenutniCvor.elements[i] != nil; i++ {
				//ako je trenutni cvoric koji gledamo veci od onog sto trazimo idemo na dete gde su elem manji
				if trenutniCvor.elements[i].key > key {
					roditelj = trenutniCvor
					trenutniCvor = *trenutniCvor.children[i]
					pronadjen = true
					break
				}
			}
			//ako ga nismo pronasli mora biti poslednje dete jer je onda kljuc veci od elem koje smo gledali
			if !pronadjen {
				roditelj = trenutniCvor
				trenutniCvor = *trenutniCvor.children[len(trenutniCvor.elements)]
			}
			pronadjen = false
		}
		//dodamo vrednost tamo gde treba
		for i := 0; i < len(trenutniCvor.elements)-1; i++ {
			if key > trenutniCvor.elements[i].key && key < trenutniCvor.elements[i+1].key {
				trenutniCvor.elements[i] = &Node{key, value}
				break
			}
		}
		//ako ima vise elem nego sto sme moramo jednog da promote-ujemo gore
		if len(trenutniCvor.elements) > 2 {
			//prvo proverimo da li ima dece
			//ako je list samo mu dodamo decu a njega smanjimo
			if len(trenutniCvor.children) == 0 {
				trenutniCvor.children = make([]*BtreeNode, 3)
				trenutniCvor.children[0] = &BtreeNode{make([]*Node, 3), make([]*BtreeNode, 3), 1}
				trenutniCvor.children[0].elements[0] = trenutniCvor.elements[0]
				trenutniCvor.children[2] = &BtreeNode{make([]*Node, 3), make([]*BtreeNode, 3), 1}
				trenutniCvor.children[2].elements[2] = trenutniCvor.elements[2]
				trenutniCvor.elements[0] = trenutniCvor.elements[1]
				trenutniCvor.elements[1] = nil
				trenutniCvor.elements[2] = nil
			} else {
				//ako ne moze tako onda treba da promote-ujemo srednji
				if roditelj.elements[0].key > trenutniCvor.elements[1].key {
					roditelj.elements[2] = roditelj.elements[1]
					roditelj.elements[1] = roditelj.elements[0]
					roditelj.elements[0] = trenutniCvor.elements[1]
				} else if roditelj.elements[1].key > trenutniCvor.elements[1].key {
					roditelj.elements[2] = roditelj.elements[1]
					roditelj.elements[1] = trenutniCvor.elements[1]
				} else {
					roditelj.elements[2] = trenutniCvor.elements[1]
				}
			}
		}
	} else {
		print("Vec postoji dati kljuc")
	}
}

func (btree *Btree) delete(key string) {
	b, elem := btree.findElement(key)
	if b == true {
	} else {
		if elem != nil {
			print("Element sa unetim kljucem je vec obrisan")
		} else {
			print("Ne postoji element sa unetim kljucem")
		}
	}

}
