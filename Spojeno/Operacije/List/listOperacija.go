package List

import (
	"sort"
	"strings"
)

func List(podstring string, velicina int, broj int, kljucevi []string) []string {
	var potrebni_kljucevi []string
	for kljuc := range kljucevi {
		if strings.HasPrefix(kljucevi[kljuc], podstring) {
			potrebni_kljucevi = append(potrebni_kljucevi, kljucevi[kljuc])
		}
	}
	sort.Strings(potrebni_kljucevi)
	var kljucevi_paginacija []string
	indeks := velicina * (broj - 1)
	for i := indeks; i < indeks+velicina; i++ {
		if i < len(potrebni_kljucevi) {
			kljucevi_paginacija = append(kljucevi_paginacija, potrebni_kljucevi[i])
		}
	}
	return kljucevi_paginacija
}
