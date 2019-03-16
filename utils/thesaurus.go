package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// MainThesaurus ..
func MainThesaurus() {

	secretKey := os.Getenv("BHT_APIKEY")

	fmt.Println(secretKey)

	thesaurus1 := &BigHugh{APIKey: secretKey}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus1.Synonyms(word)
		if err != nil {
			log.Fatalf("Failed when looking for synonyms for \""+word+"\" %v \n", err)
		}
		if len(syns) == 0 {
			log.Fatalln("Couldn't find any synonyms for \"" + word + "\"")
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}

}
