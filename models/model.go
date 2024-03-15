// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    meaning, err := UnmarshalMeaning(bytes)
//    bytes, err = meaning.Marshal()

package controllers

import "encoding/json"

type Meaning []MeaningElement

func UnmarshalMeaning(data []byte) (Meaning, error) {
	var r Meaning
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Meaning) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MeaningElement struct {
	Word       string         `json:"word"`
	Phonetic   string         `json:"phonetic"`
	Phonetics  []Phonetic     `json:"phonetics"`
	Meanings   []MeaningClass `json:"meanings"`
	License    License        `json:"license"`
	SourceUrls []string       `json:"sourceUrls"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type MeaningClass struct {
	PartOfSpeech string        `json:"partOfSpeech"`
	Definitions  []Definition  `json:"definitions"`
	Synonyms     []string      `json:"synonyms"`
	Antonyms     []interface{} `json:"antonyms"`
}

type Definition struct {
	Definition string        `json:"definition"`
	Synonyms   []interface{} `json:"synonyms"`
	Antonyms   []interface{} `json:"antonyms"`
	Example    *string       `json:"example,omitempty"`
}

type Phonetic struct {
	Text      string `json:"text"`
	Audio     string `json:"audio"`
	SourceURL string `json:"sourceUrl"`
}
