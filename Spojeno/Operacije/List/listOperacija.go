package List

import (
	"fmt"
	"strings"
)

func List(podstring string, velicina int, broj int, kljucevi []string) [][]byte {
	var potrebni_kljucevi []string
	for kljuc := range kljucevi {
		if strings.HasPrefix(kljucevi[kljuc], podstring) {
			potrebni_kljucevi = append(potrebni_kljucevi, kljucevi[kljuc])
		}
	}
	fmt.Println(potrebni_kljucevi)
	return make([][]byte, 0)
}
