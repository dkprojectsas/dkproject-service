package utils

func NumberSend(number string) (string, bool) {
	if number[0] == '0' {
		return "+62" + number[1:len(number)], true
	} else if number[0:2] == "62" {
		return "+" + number, true
	} else if number[0:3] == "+62" {
		return number, true
	} else if number[0] == '8' {
		return "+62" + number, true
	}

	return "", false
}
