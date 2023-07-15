package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type CustomNormalizer struct {
	transform.NopResetter
}

func (n *CustomNormalizer) Transform(dst, src []byte, atEOF bool) (int, int, error) {
	r, size := utf8.DecodeRune(src)
	if r == utf8.RuneError {
		return 0, 0, transform.ErrShortSrc
	}

	if r == 'ñ' || r == 'Ñ' || r == 'ü' || r == 'Ü' {
		return copy(dst, src[:size]), size, nil
	}

	normalized := norm.NFD.String(string(r))

	stripped := removeCombiningMarks(normalized)

	return copy(dst, []byte(stripped)), size, nil
}

func removeCombiningMarks(str string) string {
	var result strings.Builder
	for _, r := range str {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

func main() {
	raw := "Ilingüísticatagüíañ"
	output, _, _ := transform.String(&CustomNormalizer{}, raw)
	fmt.Println(output)
}
