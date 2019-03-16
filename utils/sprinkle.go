package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

var transforms = []string{
	"agile " + otherWord,
	otherWord + "smart",
	"wold is " + otherWord,
	otherWord + "best",
	otherWord + "app",
	otherWord + "site",
	otherWord + "time",
	"get" + otherWord,
	"go" + otherWord,
	"lets " + otherWord,
}

func mainSprinkle() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Text() == "exit" || s.Text() == "quit" {
			break
		}
		t := transforms[rand.Intn(len(transforms))]
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
}
