package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var turkish = collate.New(language.Turkish)

func TurkceCompareString(s1, s2 string) int {
	return turkish.CompareString(s1, s2)
}

// ToTitle make string to turkish title
func ToTitle(str string) string {

	dizi := strings.Split(str, " ")
	dst := ""
	for _, word := range dizi {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}

		kelime := []rune(word)
		//sadece ilk harfi büyütelim
		dst += cases.Title(language.Turkish).String(string(kelime[0]))
		dst += cases.Lower(language.Turkish).String(string(kelime[1:]))
		dst += " "
	}

	return strings.TrimSpace(dst)
}

// ToUpper make string to turkish upper
func ToUpper(str string) string {
	return cases.Upper(language.Turkish).String(strings.TrimSpace(str))
}

// ToLower make string to turkish lower
func ToLower(str string) string {
	return cases.Lower(language.Turkish).String(strings.TrimSpace(str))
}

func TurkceKarakterTemizle(str string) string {
	str = strings.ReplaceAll(str, "ç", "c")
	str = strings.ReplaceAll(str, "Ç", "C")
	str = strings.ReplaceAll(str, "ğ", "g")
	str = strings.ReplaceAll(str, "Ğ", "G")
	str = strings.ReplaceAll(str, "ı", "i")
	str = strings.ReplaceAll(str, "İ", "I")
	str = strings.ReplaceAll(str, "ö", "o")
	str = strings.ReplaceAll(str, "Ö", "O")
	str = strings.ReplaceAll(str, "ş", "s")
	str = strings.ReplaceAll(str, "Ş", "S")
	str = strings.ReplaceAll(str, "ü", "u")
	str = strings.ReplaceAll(str, "Ü", "U")

	return str
}

func EmailTemizle(email string) string {
	if email == "" {
		return ""
	}
	email = strings.TrimSpace(email)
	email = TurkceKarakterTemizle(email)
	email = strings.ToLower(email)

	return email
}

// TelefonTemizle telefon numarasını başında sıfır olmayacak şekilde temizler
func TelefonTemizle(telefon string) string {
	if telefon == "" {
		return ""
	}
	telefon = strings.TrimSpace(telefon)
	telefon = strings.ReplaceAll(telefon, " ", "")
	telefon = strings.ReplaceAll(telefon, "(", "")
	telefon = strings.ReplaceAll(telefon, ")", "")
	telefon = strings.ReplaceAll(telefon, "-", "")
	telefon = strings.ReplaceAll(telefon, "_", "")

	if strings.HasPrefix(telefon, "+90") {
		telefon = telefon[3:]
	}
	if strings.HasPrefix(telefon, "90") {
		telefon = telefon[2:]
	}

	return telefon
}

func ToSlug(sentence string) string {
	// Türkçe karakterleri uygun slug formatına dönüştürme
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	slug, _, _ := transform.String(t, sentence)

	slug = strings.Map(func(r rune) rune {
		switch r {
		case 'ı':
			return 'i'
		case 'ğ':
			return 'g'
		case 'ü':
			return 'u'
		case 'ş':
			return 's'
		case 'ö':
			return 'o'
		case 'ç':
			return 'c'
		default:
			return r
		}
	}, slug)

	// Diğer tüm karakterleri çıkarılması
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, ".", "-")

	// Slug'un başındaki ve sonundaki tireleri kaldırma
	slug = strings.Trim(slug, "-")

	return slug
}
