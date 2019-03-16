package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Thesaurus ...
type Thesaurus interface {
	Synonyms(term string) ([]string, error)
}

// BigHugh ...
type BigHugh struct {
	APIKey string
}
type synonyms struct {
	Noun      *words `json:"noun"`
	Verb      *words `json:"verb"`
	Adjective *words `json:"adjective"`
}
type words struct {
	Syn []string `json:"syn"`
}

// Synonyms ...
func (b *BigHugh) Synonyms(term string) ([]string, error) {
	log.Println("Looking for term ", term)
	var syns []string
	response, err := http.Get("http://words.bighugelabs.com/api/2/" + b.APIKey + "/" + term + "/json")
	if err != nil {
		return syns, errors.New("bighugh: Failed when looking for synonyms for \"" + term + "\"" + err.Error())
	}
	var data synonyms

	// if response.StatusCode != http.StatusOK {
	// 	bodyBytes, _ := ioutil.ReadAll(response.Body)
	// 	bodyString := string(bodyBytes)
	// 	log.Println("Looking for term body is ", bodyString)
	// } else {
	// 	log.Println("Status code not ok ", response.StatusCode)
	// }

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return syns, errors.New(strconv.Itoa(response.StatusCode))
	}
	if err := json.NewDecoder(response.Body).Decode(&data); err !=
		nil {
		return syns, err
	}
	if data.Noun != nil {
		syns = append(syns, data.Noun.Syn...)
	}
	if data.Verb != nil {
		syns = append(syns, data.Verb.Syn...)
	}
	if data.Adjective != nil{
		syns = append(syns, data.Adjective.Syn...)
	}
	return syns, nil
}
