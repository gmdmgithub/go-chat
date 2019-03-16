package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var tlds = []string{"com", "net", "org"}

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"

func mainDomainify() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := strings.ToLower(s.Text())
		var newText []rune
		var isPreviousSpace = false
		for index, r := range text {
			if unicode.IsSpace(r) {
				if isPreviousSpace {
					continue
				}
				isPreviousSpace = true
				if index != 0 && index != (len(text)-1) {
					r = '-'
				}
			} else {
				isPreviousSpace = false
			}
			if !strings.ContainsRune(allowedChars, r) {
				continue
			}
			newText = append(newText, r)
		}
		fmt.Println(string(newText) + "." +
			tlds[rand.Intn(len(tlds))])
	}
}
