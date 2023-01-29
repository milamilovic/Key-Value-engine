package main

import "fmt"

// func main() {
// 	btree := makeBtree(3)
// 	btree.add("3", []byte("a"))
// 	fmt.Println(btree.root)
// 	fmt.Println("ovde sam")
// 	btree.add("vau", []byte("a"))
// 	fmt.Println("ovde sam 2")
// 	btree.add("ne znam", []byte("a"))
// 	fmt.Println(btree.root)
// 	btree.add("kuku", []byte("a"))
// 	btree.add("fff", []byte("a"))
// 	btree.add("1", []byte("ponovo isti kljuc"))
// 	btree.add("4", []byte("a"))
// 	fmt.Println()
// 	fmt.Println()
// 	fmt.Println(btree.root)
// 	fmt.Println(btree.root.children[0])
// 	fmt.Println(btree.root.children[1])
// 	fmt.Println()
// 	fmt.Println()
// 	b, _ := btree.findElement("3")
// 	print(b)
// 	fmt.Println()
// 	fmt.Println()
// 	b, _ = btree.findElement("vau")
// 	print(b)
// 	fmt.Println()
// 	fmt.Println()
// 	b, _ = btree.findElement("fff")
// 	print(b)
// 	fmt.Println()
// 	fmt.Println()
// 	b, _ = btree.findElement("12") //ovo je false
// 	print(b)
// 	fmt.Println()
// 	fmt.Println()
// 	b, _ = btree.findElement("4")
// 	print(b)
// }

type Btree struct {
	root      *BtreeNode
	maxHeight int
	height    int
}

type BtreeNode struct {
	elements     []*Node
	children     []*BtreeNode
	max_children int
	num_of_elem  int
	parent       *BtreeNode
}

type Node struct {
	key   string
	value []byte
}

func makeBtree(maxHeight int) *Btree {
	i := 0
	j := 0
	k := 1
	root := BtreeNode{make([]*Node, 4), make([]*BtreeNode, 4), i, j, nil}
	return &Btree{&root, maxHeight + 1, k}
}

func (btree *Btree) findElement(key string) (bool, *Node) {
	//fmt.Println("trazimo elem sa kljucem ", key)
	trenutniCvor := *btree.root
	if trenutniCvor.max_children == 0 {
		return false, nil
	}
	k := 0
	for i := 0; i <= btree.height; i++ {
		trenutnoDete := trenutniCvor.elements[k]
		//fmt.Println("kljuc trenutnog deteta je ", trenutnoDete.key)
		sledeceDete := trenutniCvor.elements[k+1]
		for z := 0; z < trenutniCvor.num_of_elem; z++ {
			trenutnoDete = trenutniCvor.elements[k]
			//fmt.Println("kljuc trenutnog deteta je ", trenutnoDete.key)
			sledeceDete = trenutniCvor.elements[k+1]
			//ako je ovaj dete tjt nasli smo
			//fmt.Println("kljuc koji trazimo je ", key)
			//fmt.Println(key == trenutnoDete.key)
			if trenutnoDete.key == key {
				//fmt.Println("nasli smo trenutno dete")
				return true, trenutnoDete
			}
			if sledeceDete == nil {
				break
			}
			//fmt.Println("kljuc sledeceg deteta je ", sledeceDete.key)
			//ako je sledeci dete tjt nasli smo
			if sledeceDete.key == key {
				//fmt.Println("nasli smo sledece dete")
				return true, sledeceDete
			}
			//ako je kljuc manji od trenutnog kljuca idemo na dete tog indeksa
			if key < trenutnoDete.key {
				//fmt.Println("kljuc je manji od kljuca trenutnog deteta")
				if trenutniCvor.children[k] == nil {
					//fmt.Println("posto trenutni cvor nema dece znaci da nema naseg elem")
					//fmt.Println(trenutniCvor.children)
					//fmt.Println("k je ", k)
					return false, nil
				}
				//fmt.Println("Trenutni cvor postaje dete broj ", k, " tj ", trenutniCvor.children[k])
				trenutniCvor = *trenutniCvor.children[k]
				k = 0
				continue
			}
			//ako je kljuc izmedju trenutnog i sledeceg idemo tamo
			if sledeceDete.key > key && trenutnoDete.key < key {
				//fmt.Println("kljuc je izmedju trenutnog i sledeceg deteta")
				if trenutniCvor.children[k+1] == nil {
					//fmt.Println("posto nema tog deteta necemo nista da radimo")
					return false, nil
				}
				trenutniCvor = *trenutniCvor.children[k+1]
				k = 0
				continue
			}
			k++
		}
		//ako je sl dete nil onda ili nema deteta ili je ,,desno"
		//fmt.Println("nema narednog deteta")
		if key > trenutnoDete.key {
			//fmt.Println("kljhuc je veci od kljuca trenutnog deteta koje je poslednje")
			if k < 3 {
				if trenutniCvor.children[3] != nil {
					//fmt.Println("trenutni cvor postaje cetvrto dete")
					trenutniCvor = *trenutniCvor.children[3]
				} else if trenutniCvor.children[2] != nil {
					//fmt.Println("trenutni cvor postaje trece dete")
					trenutniCvor = *trenutniCvor.children[2]
				} else if trenutniCvor.children[1] != nil {
					//fmt.Println("trenutni cvor postaje drugo dete")
					trenutniCvor = *trenutniCvor.children[1]
				} else if trenutniCvor.children[0] != nil {
					//fmt.Println("trenutni cvor postaje rpvo dete")
					trenutniCvor = *trenutniCvor.children[0]
				} else {
					return false, nil
				}
			}
			k = 0
		}
	}
	return false, nil //nismo upste nasli
}

func (btree *Btree) add(key string, value []byte) {
	b, _ := btree.findElement(key)
	trenutniCvor := *btree.root
	pokaz_na_trenutni_cvor := btree.root
	//roditelj := trenutniCvor
	red_tr_cvora_kod_roditelja := 0
	pokaz_na_roditelja := btree.root
	// fmt.Println(pokaz_na_roditelja)
	// fmt.Println("da li elem vec postoji: ", b)
	// fmt.Println("dodajemo element, trenutno stanje korena")
	// fmt.Println(trenutniCvor)
	if b == false { //ovo znaci da ne postoji jer ga nije pronasao
		//fmt.Println("pronadjen element")
		pronadjen := false
		//fmt.Println(trenutniCvor.children)
		//idemo sve do lista
		for trenutniCvor.children[0] != nil {
			//fmt.Println(trenutniCvor)
			//fmt.Println("ovo nije list")
			for i := 0; i < trenutniCvor.num_of_elem-1 && trenutniCvor.elements[i] != nil; i++ {
				//ako je trenutni cvoric koji gledamo veci od onog sto trazimo idemo na dete gde su elem manji
				//fmt.Println("usao u for")
				if trenutniCvor.elements[i].key > key {
					//fmt.Println("usao u if")
					//roditelj = trenutniCvor
					pokaz_na_roditelja = pokaz_na_trenutni_cvor
					trenutniCvor = *trenutniCvor.children[i]
					pokaz_na_trenutni_cvor = trenutniCvor.children[i]
					red_tr_cvora_kod_roditelja = i
					pronadjen = true
					break
				}
			}
			//ako ga nismo pronasli mora biti poslednje dete jer je onda kljuc veci od elem koje smo gledali
			if !pronadjen {
				//roditelj = trenutniCvor
				if trenutniCvor.children[trenutniCvor.num_of_elem] != nil {
					trenutniCvor = *trenutniCvor.children[trenutniCvor.num_of_elem]
					pokaz_na_trenutni_cvor = trenutniCvor.children[trenutniCvor.num_of_elem]
					pokaz_na_roditelja = pokaz_na_trenutni_cvor
					red_tr_cvora_kod_roditelja = trenutniCvor.num_of_elem
				} else {
					trenutniCvor = *trenutniCvor.children[0]
					pokaz_na_trenutni_cvor = trenutniCvor.children[0]
					pokaz_na_roditelja = pokaz_na_trenutni_cvor
					red_tr_cvora_kod_roditelja = trenutniCvor.num_of_elem
				}
			}
			pronadjen = false
		}
		//dodamo vrednost tamo gde treba
		//fmt.Println("na listu smo")
		//ako cvor nema elem tu dodajemo
		if trenutniCvor.num_of_elem == 0 {
			// fmt.Println("nema elem na cvoru na kom smo")
			// fmt.Println("stanje pre dodavanja")
			// fmt.Println(trenutniCvor)
			trenutniCvor.elements[0] = &Node{key, value}
			trenutniCvor.num_of_elem = trenutniCvor.num_of_elem + 1
			trenutniCvor.max_children = trenutniCvor.num_of_elem + 1
			if pokaz_na_roditelja == pokaz_na_trenutni_cvor {
				//onda smo na korenu treba da izmenimo koren
				pokaz_na_trenutni_cvor = &trenutniCvor
				pokaz_na_roditelja = pokaz_na_trenutni_cvor
				btree.root = pokaz_na_trenutni_cvor
			} else {
				//inace menjamo kod roditelja pokazivac na dete koje gledamo
				pokaz_na_trenutni_cvor = &trenutniCvor
				pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
			}
			// fmt.Println("stanje posle dodavanja, cvor")
			// fmt.Println(trenutniCvor)
			// fmt.Println("stanje posle dodavanja, pokaz na cvor")
			// fmt.Println(*pokaz_na_trenutni_cvor)
			// fmt.Println("stanje posle dodavanja, koren")
			// fmt.Println(btree.root)
		} else {
			// fmt.Println("ima bar 1 elem")
			// fmt.Println(trenutniCvor.elements)
			//iteriramo
			menjao_se_u_petlji := false
			for i := 0; i < 3; i++ {
				//fmt.Println("i je ", i)
				if trenutniCvor.elements[i+1] != nil {
					// fmt.Println("usli smo tamo gde se iterira za dodavanje")
					// fmt.Println(key)
					// fmt.Println(trenutniCvor.elements[i].key)
					// fmt.Println(trenutniCvor.elements[i+1].key)
					if key > trenutniCvor.elements[i].key && key < trenutniCvor.elements[i+1].key {
						// fmt.Println("cao")
						// fmt.Println("stanje pre dodavanja")
						// fmt.Println(trenutniCvor)
						for m := 3; m > i; m-- {
							trenutniCvor.elements[m] = trenutniCvor.elements[m-1]
							//fmt.Println(trenutniCvor)
						}
						trenutniCvor.elements[i+1] = &Node{key, value}
						trenutniCvor.num_of_elem = trenutniCvor.num_of_elem + 1
						trenutniCvor.max_children = trenutniCvor.num_of_elem + 1
						if pokaz_na_roditelja == pokaz_na_trenutni_cvor {
							//onda smo na koreno
							pokaz_na_trenutni_cvor = &trenutniCvor
							pokaz_na_roditelja = pokaz_na_trenutni_cvor
							btree.root = pokaz_na_trenutni_cvor
						} else {
							pokaz_na_trenutni_cvor = &trenutniCvor
							pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
						}
						// fmt.Println("stanje posle dodavanja, cvor")
						// fmt.Println(trenutniCvor)
						// fmt.Println("stanje posle dodavanja, pokaz na cvor")
						// fmt.Println(*pokaz_na_trenutni_cvor)
						menjao_se_u_petlji = true
						break
					}

				} else if trenutniCvor.num_of_elem == 1 {
					//imamo jedan element i treba da dodamo drugi
					// fmt.Println("imamo 1 elem i dodajemo drugi")
					// fmt.Println(trenutniCvor)
					if key > trenutniCvor.elements[0].key {
						//ako je veci od prethodnog dodajemo ga posle
						// fmt.Println("stanje pre dodavanja")
						// fmt.Println(trenutniCvor)
						trenutniCvor.elements[1] = &Node{key, value}
						trenutniCvor.num_of_elem = trenutniCvor.num_of_elem + 1
						trenutniCvor.max_children = trenutniCvor.num_of_elem + 1
						if pokaz_na_roditelja == pokaz_na_trenutni_cvor {
							//onda smo na koreno
							pokaz_na_trenutni_cvor = &trenutniCvor
							pokaz_na_roditelja = pokaz_na_trenutni_cvor
							btree.root = pokaz_na_trenutni_cvor
						} else {
							pokaz_na_trenutni_cvor = &trenutniCvor
							pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
						}
						// fmt.Println("stanje posle dodavanja, cvor")
						// fmt.Println(trenutniCvor)
						// fmt.Println("stanje posle dodavanja, pokaz na cvor")
						// fmt.Println(*pokaz_na_trenutni_cvor)
						menjao_se_u_petlji = true
						break
					} else {
						trenutniCvor.elements[1] = trenutniCvor.elements[0]
						trenutniCvor.elements[0] = &Node{key, value}
						trenutniCvor.num_of_elem = trenutniCvor.num_of_elem + 1
						trenutniCvor.max_children = trenutniCvor.num_of_elem + 1
						if pokaz_na_roditelja == pokaz_na_trenutni_cvor {
							//onda smo na koreno
							pokaz_na_trenutni_cvor = &trenutniCvor
							pokaz_na_roditelja = pokaz_na_trenutni_cvor
							btree.root = pokaz_na_trenutni_cvor
						} else {
							pokaz_na_trenutni_cvor = &trenutniCvor
							pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
						}
						// fmt.Println("stanje posle dodavanja, cvor")
						// fmt.Println(trenutniCvor)
						// fmt.Println("stanje posle dodavanja, pokaz na cvor")
						// fmt.Println(*pokaz_na_trenutni_cvor)
						menjao_se_u_petlji = true
						break
					}
				} else if (i == 2) && trenutniCvor.num_of_elem == 2 {
					// fmt.Println("treba da dodamo na kraj")
					// fmt.Println("stanje pre dodavanja")
					// fmt.Println(trenutniCvor)
					if trenutniCvor.elements[i] == nil {
						//fmt.Println("lalala tu smo usli")
						trenutniCvor.elements[i] = &Node{key, value}
						trenutniCvor.num_of_elem = trenutniCvor.num_of_elem + 1
						trenutniCvor.max_children = trenutniCvor.num_of_elem + 1

						if pokaz_na_roditelja == pokaz_na_trenutni_cvor {
							//onda smo na koreno
							pokaz_na_trenutni_cvor = &trenutniCvor
							btree.root = pokaz_na_trenutni_cvor
						} else {
							pokaz_na_trenutni_cvor = &trenutniCvor
							pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
						}
						// fmt.Println("stanje posle dodavanja, cvor")
						// fmt.Println(trenutniCvor)
						// fmt.Println("stanje posle dodavanja, pokaz na cvor")
						// fmt.Println(*pokaz_na_trenutni_cvor)
						menjao_se_u_petlji = true
						break
					}
				}
			}
			if trenutniCvor.num_of_elem == 3 && !menjao_se_u_petlji {
				trenutniCvor.elements[3] = &Node{key, value}
				trenutniCvor.num_of_elem = trenutniCvor.num_of_elem + 1
				trenutniCvor.max_children = trenutniCvor.num_of_elem + 1

				if pokaz_na_roditelja == pokaz_na_trenutni_cvor {
					//onda smo na koreno
					pokaz_na_trenutni_cvor = &trenutniCvor
					pokaz_na_roditelja = &trenutniCvor
					btree.root = pokaz_na_trenutni_cvor
				} else {
					pokaz_na_trenutni_cvor = &trenutniCvor
					pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
				}
				// fmt.Println("stanje posle dodavanja, cvor")
				// fmt.Println(trenutniCvor)
				// fmt.Println("stanje posle dodavanja, pokaz na cvor")
				// fmt.Println(*pokaz_na_trenutni_cvor)
			}
		}
		//ako ima vise elem nego sto sme moramo jednog da promote-ujemo gore
		//fmt.Println(trenutniCvor.num_of_elem)
		if trenutniCvor.num_of_elem > 3 {
			//prvo proverimo da li ima dece
			//ako je list samo mu dodamo decu a njega smanjimo
			// fmt.Println("treba da promote-ujemo nekoga, trenutni cvor: ")
			// fmt.Println(trenutniCvor)
			// fmt.Println("deca trenutnog: ")
			// fmt.Println(trenutniCvor.children)

			//NE RADI OVO DELJENJE
			//NE MENJA DOBRO, RADI KAO DA OD 3 PRAVI 2 A TREBA OD 4 3 VALJA ISTRAZITI!!!

			//ako list treba da se deli onda samo ostavimo srednja 2 gore a ostale smestimo u decu levo i desno
			if trenutniCvor.children[0] == nil {
				trenutniCvor.children[0] = &BtreeNode{make([]*Node, 4), make([]*BtreeNode, 4), 0, 1, &trenutniCvor}
				trenutniCvor.children[0].elements[0] = trenutniCvor.elements[0]
				trenutniCvor.children[1] = &BtreeNode{make([]*Node, 4), make([]*BtreeNode, 4), 0, 1, &trenutniCvor}
				trenutniCvor.children[1].elements[0] = trenutniCvor.elements[3]
				fmt.Println(trenutniCvor)
				trenutniCvor.elements[0] = trenutniCvor.elements[1]
				trenutniCvor.elements[1] = trenutniCvor.elements[2]
				trenutniCvor.elements[2] = nil
				trenutniCvor.elements[3] = nil
				trenutniCvor.num_of_elem = 2
				trenutniCvor.max_children = 3
				// fmt.Println(pokaz_na_roditelja)
				// fmt.Println(pokaz_na_trenutni_cvor)
				if pokaz_na_roditelja == pokaz_na_trenutni_cvor {
					//onda smo na koreno
					//fmt.Println("na korenu smo")
					pokaz_na_trenutni_cvor = &trenutniCvor
					pokaz_na_roditelja = &trenutniCvor
					btree.root = &trenutniCvor
				} else {
					pokaz_na_trenutni_cvor = &trenutniCvor
					pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
				}
				// fmt.Println("stanje posle dodavanja, cvor")
				// fmt.Println(trenutniCvor)
				// fmt.Println("stanje posle dodavanja, pokaz na cvor")
				// fmt.Println(*pokaz_na_trenutni_cvor)
				// fmt.Println("koren")
				// fmt.Println(btree.root)
				// fmt.Println("stanje posle dodavanja, dete 1")
				// fmt.Println(trenutniCvor.children[0])
				// fmt.Println("stanje posle dodavanja, dete 2")
				// fmt.Println(trenutniCvor.children[1])
			} else {
				//ako ne moze tako onda treba da promote-ujemo srednji
				for trenutniCvor.num_of_elem != 3 {
					broj_dece := 0
					for r := 0; r < 4; r++ {
						if trenutniCvor.children[r] != nil {
							broj_dece++
						}
					}
					//ako je poslednje dete nepopunjeno tamo dodamo jedan elem iz roditelja
					if trenutniCvor.children[broj_dece-1].num_of_elem < 3 {
						broj_dece2 := 0
						for r := 0; r < 4; r++ {
							if trenutniCvor.children[broj_dece-1].elements[r] != nil {
								broj_dece2++
							}
						}
						trenutniCvor.children[broj_dece-1].elements[broj_dece2-1] = trenutniCvor.elements[3]
						trenutniCvor.children[broj_dece-1].num_of_elem++
						trenutniCvor.children[broj_dece-1].max_children = trenutniCvor.children[broj_dece-1].num_of_elem + 1
						pokaz_na_trenutni_cvor = &trenutniCvor
						pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
					}
					//ako je prvo dete nepopunjeno tamo dodajemo elem iz roditelja
					if broj_dece != 1 {
						if trenutniCvor.children[0].num_of_elem < 3 {
							broj_dece3 := 0
							for r := 0; r < 4; r++ {
								if trenutniCvor.children[0].elements[r] != nil {
									broj_dece3++
								}
							}
							trenutniCvor.children[0].elements[broj_dece3-1] = trenutniCvor.elements[3]
							trenutniCvor.children[0].num_of_elem++
							trenutniCvor.children[0].max_children = trenutniCvor.children[0].num_of_elem + 1
							pokaz_na_trenutni_cvor = &trenutniCvor
							pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
						}
					}
					//ako ne moze ni to onda dodajemo u roditelja i menjamo koji je trenutni cvor
					//pa se while petlja nastavlja
					if red_tr_cvora_kod_roditelja == 0 {
						pokaz_na_roditelja.elements[3] = pokaz_na_roditelja.elements[2]
						pokaz_na_roditelja.elements[2] = pokaz_na_roditelja.elements[1]
						pokaz_na_roditelja.elements[1] = pokaz_na_roditelja.elements[0]
						pokaz_na_roditelja.elements[0] = trenutniCvor.elements[3]
						trenutniCvor.elements[3] = nil
						pokaz_na_roditelja.num_of_elem++
						pokaz_na_roditelja.max_children = pokaz_na_roditelja.num_of_elem + 1
						trenutniCvor.num_of_elem = 3
						trenutniCvor.max_children = 4
					} else if red_tr_cvora_kod_roditelja == 3 {
						pokaz_na_roditelja.elements[3] = trenutniCvor.elements[0]
						trenutniCvor.elements[0] = trenutniCvor.elements[1]
						trenutniCvor.elements[1] = trenutniCvor.elements[2]
						trenutniCvor.elements[2] = trenutniCvor.elements[3]
						trenutniCvor.elements[3] = nil
						pokaz_na_roditelja.num_of_elem++
						pokaz_na_roditelja.max_children = pokaz_na_roditelja.num_of_elem + 1
						trenutniCvor.num_of_elem = 3
						trenutniCvor.max_children = 4
					} else if red_tr_cvora_kod_roditelja == 2 {
						if pokaz_na_roditelja.num_of_elem == 2 {
							pokaz_na_roditelja.elements[2] = trenutniCvor.elements[0]
							trenutniCvor.elements[0] = trenutniCvor.elements[1]
							trenutniCvor.elements[1] = trenutniCvor.elements[2]
							trenutniCvor.elements[2] = trenutniCvor.elements[3]
							trenutniCvor.elements[3] = nil
							pokaz_na_roditelja.num_of_elem++
							pokaz_na_roditelja.max_children = pokaz_na_roditelja.num_of_elem + 1
							trenutniCvor.num_of_elem = 3
							trenutniCvor.max_children = 4
						} else if pokaz_na_roditelja.num_of_elem == 3 {
							pokaz_na_roditelja.elements[3] = pokaz_na_roditelja.elements[2]
							pokaz_na_roditelja.elements[2] = trenutniCvor.elements[0]
							trenutniCvor.elements[0] = trenutniCvor.elements[1]
							trenutniCvor.elements[1] = trenutniCvor.elements[2]
							trenutniCvor.elements[2] = trenutniCvor.elements[3]
							trenutniCvor.elements[3] = nil
							pokaz_na_roditelja.num_of_elem++
							pokaz_na_roditelja.max_children = pokaz_na_roditelja.num_of_elem + 1
							trenutniCvor.num_of_elem = 3
							trenutniCvor.max_children = 4
						}
					} else if red_tr_cvora_kod_roditelja == 1 {
						if pokaz_na_roditelja.num_of_elem == 1 {
							pokaz_na_roditelja.elements[1] = trenutniCvor.elements[0]
							trenutniCvor.elements[0] = trenutniCvor.elements[1]
							trenutniCvor.elements[1] = trenutniCvor.elements[2]
							trenutniCvor.elements[2] = trenutniCvor.elements[3]
							trenutniCvor.elements[3] = nil
							pokaz_na_roditelja.num_of_elem++
							pokaz_na_roditelja.max_children = pokaz_na_roditelja.num_of_elem + 1
							trenutniCvor.num_of_elem = 3
							trenutniCvor.max_children = 4
						} else if pokaz_na_roditelja.num_of_elem == 2 {
							pokaz_na_roditelja.elements[2] = pokaz_na_roditelja.elements[1]
							pokaz_na_roditelja.elements[1] = trenutniCvor.elements[0]
							trenutniCvor.elements[0] = trenutniCvor.elements[1]
							trenutniCvor.elements[1] = trenutniCvor.elements[2]
							trenutniCvor.elements[2] = trenutniCvor.elements[3]
							trenutniCvor.elements[3] = nil
							pokaz_na_roditelja.num_of_elem++
							pokaz_na_roditelja.max_children = pokaz_na_roditelja.num_of_elem + 1
							trenutniCvor.num_of_elem = 3
							trenutniCvor.max_children = 4
						} else if pokaz_na_roditelja.num_of_elem == 3 {
							pokaz_na_roditelja.elements[3] = pokaz_na_roditelja.elements[2]
							pokaz_na_roditelja.elements[2] = pokaz_na_roditelja.elements[1]
							pokaz_na_roditelja.elements[1] = trenutniCvor.elements[0]
							trenutniCvor.elements[0] = trenutniCvor.elements[1]
							trenutniCvor.elements[1] = trenutniCvor.elements[2]
							trenutniCvor.elements[2] = trenutniCvor.elements[3]
							trenutniCvor.elements[3] = nil
							pokaz_na_roditelja.num_of_elem++
							pokaz_na_roditelja.max_children = pokaz_na_roditelja.num_of_elem + 1
							trenutniCvor.num_of_elem = 3
							trenutniCvor.max_children = 4
						}
					}
					pokaz_na_trenutni_cvor = trenutniCvor.children[0]
					pokaz_na_roditelja = pokaz_na_trenutni_cvor.parent
					pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = pokaz_na_trenutni_cvor
					if pokaz_na_roditelja == pokaz_na_trenutni_cvor {
						//onda smo na koreno
						//fmt.Println("na korenu smo")
						btree.root = pokaz_na_roditelja
					}
					// pokaz_na_roditelja.children[red_tr_cvora_kod_roditelja] = &trenutniCvor
					// pokaz_na_trenutni_cvor = pokaz_na_roditelja
					// pokaz_na_roditelja = pokaz_na_trenutni_cvor.parent
				}
			}
		}
	} else {
		fmt.Println("Vec postoji dati kljuc")
	}
}
