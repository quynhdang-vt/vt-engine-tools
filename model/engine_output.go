package model

import (
	"encoding/json"
	"io/ioutil"
	"sort"
)

/**
see https://github.com/veritone/database/blob/master/engine-output/sample.js#L192
*/
const ENGINE_OUTPUT_SCHEMAID = "https://www.veritone.com/schema/engine/20180524"
const MYENGINE_GUID = "ad7d9d36-dd10-4448-ac79-63b4817bd3e6"

type UtteranceWord struct {
	Word            string  `json:"word"`
	Confidence      float32 `json:"confidence"`
	BestPath        bool    `json:"bestPath"`
	UtteranceLength int     `json:"utteranceLength"`
}
type Utterance struct {
	StartTimeMs int             `json:"startTimeMs"`
	StopTimeMs  int             `json:"stopTimeMs"`
	SpeakerId   string          `json:"speakerId"`
	Words       []UtteranceWord `json:"words"`
}

type TranscriptionEngineOutput struct {
	SchemaId       string      `json:"schemaId"`
	SourceEngineId string      `json:"sourceEngineId"`
	Series         []Utterance `json:"series"`
}

func NewEngineOutput() TranscriptionEngineOutput {
	return TranscriptionEngineOutput{
		SchemaId:       ENGINE_OUTPUT_SCHEMAID,
		SourceEngineId: MYENGINE_GUID,
		Series:         make([]Utterance, 0, 10),
	}
}

func (this *TranscriptionEngineOutput) addVLFUtterance(speakerId string, vlfUtterance VLFUtterance) {
	nWords := len(vlfUtterance.Words)
	if nWords == 0 {
		return
	}
	words := make([]UtteranceWord, nWords)
	for i, v := range vlfUtterance.Words {
		words[i].Word = v.Word
		words[i].Confidence = float32(v.Confidence) / 1000.0
		words[i].BestPath = v.BestPathForward
		words[i].UtteranceLength = v.SpanningLength
	}
	utterance := Utterance{
		SpeakerId:   speakerId,
		StartTimeMs: vlfUtterance.StartTimeMs,
		StopTimeMs:  vlfUtterance.StopTimeMs,
		Words:       words,
	}
	this.Series = append(this.Series, utterance)
}

func (this *TranscriptionEngineOutput) AddVLF(speakerId string, lattice VLFLattice) {
	for _, v := range lattice {
		this.addVLFUtterance(speakerId, v)
	}
}

// sorting based on startTimeMs
func (t TranscriptionEngineOutput) Len() int {
	return len(t.Series)
}

func (t TranscriptionEngineOutput) Less(i, j int) bool {
	return t.Series[i].StartTimeMs < t.Series[j].StartTimeMs
}
func (t TranscriptionEngineOutput) Swap(i, j int) {
	t.Series[i], t.Series[j] = t.Series[j], t.Series[i]
}
func (t TranscriptionEngineOutput) Sort() {
	sort.Sort(t)
}
func (t TranscriptionEngineOutput) Write(filename string) error {
	output, _ := json.Marshal(t)
	return ioutil.WriteFile(filename, output, 0644)
}
