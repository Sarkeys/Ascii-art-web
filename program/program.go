package program

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"regexp"
)

func HashChecker(font string) error {
	font = "program/banners/" + font
	hasher := sha256.New()
	data, err := os.ReadFile(font)
	if err != nil {
		fmt.Println(err)
		return err
	}
	hasher.Write(data)
	generatedHash := fmt.Sprintf("%x", hasher.Sum(nil))
	hashMap := map[string]string{
		"program/banners/standard.txt":   "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf",
		"program/banners/shadow.txt":     "26b94d0b134b77e9fd23e0360bfd81740f80fb7f6541d1d8c5d85e73ee550f73",
		"program/banners/thinkertoy.txt": "64285e4960d199f4819323c4dc6319ba34f1f0dd9da14d07111345f5d76c3fa3",
	}
	
	if hashMap[font] != generatedHash {
		return fmt.Errorf("Invalid %s font.", font[8:len(font)-4])
	}
	return nil
}

func Convert(input, fontName string) (string, error) {
	err := HashChecker(fontName)
	if err != nil {
		return "", err
	}

	text := NewLineBreaker(input)

	font, err := GetFont("program/banners/" + fontName)
	if err != nil {
		return "", err
	}

	result := AsciiArt(text, font)
	return result, nil
}

func GetFont(fontName string) (map[rune][]string, error) {
	file, err := os.Open(fontName)
	if err != nil {
		return map[rune][]string{}, err
	}

	count := 0
	r := ' '
	font := make(map[rune][]string, 8)

	scanner := bufio.NewScanner(file)
	for ; scanner.Scan(); count++ {
		text := scanner.Text()
		if count != 0 && count != 9 {
			font[r] = append(font[r], text)
		} else if count == 9 {
			count = 0
			r++
		}
	}
	return font, nil
}

func AsciiArt(text []string, font map[rune][]string) string {
	result := ""
	output := make([][8]string, len(text))
	for k, val := range text {
		if val == "" {
			output[k] = [8]string{"\n"}
			continue
		}
		for _, r := range val {
			for i, g := range font[r] {
				output[k][i] += g
			}
		}
	}

	for _, h := range output {
		for _, k := range h {
			if k == "\n" {
				result += "\n"
				break
			}
			result += k + "\n"
		}
	}

	return result
}

func NewLineBreaker(input string) []string {
	if len(input) == 0 {
		return []string{}
	}

	input = regexp.MustCompile(`\\\\`).ReplaceAllString(input, "♥")
	input = regexp.MustCompile(`\\n`).ReplaceAllString(input, "\n")
	input = regexp.MustCompile(`\\"`).ReplaceAllString(input, "\"")
	text := []string{}
	word := ""
	count := 0

	for _, r := range input {
		if r == '\n' {
			count++
		}
	}

	if count == len(input) {
		return []string{}
	}

	for i, r := range input {
		if r == '\n' {
			text = append(text, word)
			word = ""

			if i == len(input)-1 {
				text = append(text, word)
			}

			continue
		}

		word += string(r)
	}

	if len(word) != 0 {
		text = append(text, word)
	}

	for i := 0; i < len(text); i++ {
		text[i] = regexp.MustCompile(`♥`).ReplaceAllString(text[i], "\\")
	}

	return text
}
