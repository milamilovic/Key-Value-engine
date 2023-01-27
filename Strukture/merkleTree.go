package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

func main() {
	
}

type Hash [20]byte    //vrednost

// hash funkcije
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func hash(podaci []byte) Hash {
	return sha1.Sum(podaci)
}

type Node struct {
	levi  Hashable
	desni Hashable
}

type Hashable struct {
	hash() Hash
}

// hesiranje blokova sa podacima
type Blok string

func (b Blok) hash() Hash {
	return hash([]byte(b)[:])
}

// prazan blok se pravi ako je ukupan broj blokova neparan (onda jedan u paru (desni) mora biti prazan)
type PrazanBlok struct {
}

func (pb PrazanBlok) hash() Hash {
	return [20]byte{}
}

func (n Node) hash() Hash {
	var l, d [sha1.Size]byte
	l = n.levi.hash()
	d = n.desni.hash()
	return hash(append(l[:], d[:]...))
}

// funkcija za kreiranje merkl stabla
func kreirajMerkleTree(delovi []Hashable) []Hashable {
	var nodes []Hashable
	var i int
	for i = 0; i < len(delovi); i += 2 {
		if i + 1 < len(delovi) {
			nodes = append(nodes, Node{levi: delovi[i], desni: delovi[i+1]})
		} else {
			nodes = append(nodes, Node{levi: delovi[i], desni: PrazanBlok{}})
		}
	}
	// koren
	if len(nodes) == 1 {
		return nodes
	// ostali cvorovi
	} else if len(nodes) > 1 {
		return kreirajMerkleTree(nodes)
	}
}