package validator

import (
	"slices"
	"strings"
	"unicode"
)

func NormolizePhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)

	var nomolizedNumber = ""
	var i, g int
	if strings.HasPrefix(phone, "+7") {
		nomolizedNumber += "+7"
		i += 2
		g += 1
	} else if strings.HasPrefix(phone, "8") {
		nomolizedNumber += "8"
		i += 1
		g += 1
	} else {
		return ""
	}

	for i < len(phone) {
		if slices.Contains([]string{"(", ")", "-", " "}, string(phone[i])) {
			i++
			continue
		}
		if slices.Contains([]int{1, 4, 7, 9}, g) {
			nomolizedNumber += " "
		}
		if unicode.IsDigit(rune(phone[i])) {
			nomolizedNumber += string(phone[i])
			i++
			g++
		}
	}
	return nomolizedNumber
}
