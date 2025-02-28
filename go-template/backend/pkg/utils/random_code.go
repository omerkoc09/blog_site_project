package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomNumber(length int) string {
	kod := ""
	for i := 0; i < length; i++ {
		r, _ := rand.Int(rand.Reader, big.NewInt(10))
		kod += r.String()
	}
	return kod
}

func GenerateRandomChracter(length int) string {
	harfler := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" //"abcdefghijklmnopqrstuvwxyz"
	hlen := len(harfler)
	rasgele := ""
	for i := 0; i < length; i++ {
		if i > 0 && i%2 == 0 {
			// araya rakam koyalım ki küfür vs çıkmasın
			r, _ := rand.Int(rand.Reader, big.NewInt(10))
			rasgele += r.String()
			continue
		}
		rasgele += string(harfler[getNumber(hlen)])
	}
	return rasgele
}

func getNumber(max int) int {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(r.Int64())
}
