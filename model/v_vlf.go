package model

import (
	"encoding/json"
	"io/ioutil"
)

type VLFWord struct {
	Word             string `json:"word"`
	Confidence       int    `json:"confidence"`
	BestPathForward  bool   `json:"bestPathForward"`
	BestPathBackward bool   `json:"bestPathBackward"`
	SpanningForward  bool   `json:"spanningForward"`
	SpanningBackward bool   `json:"spanningBackward"`
	SpanningLength   int    `json:"spanningLength"`
}

type VLFUtterance struct {
	Index       int       `json:"index"`
	StartTimeMs int       `json:"startTimeMs"`
	StopTimeMs  int       `json:"stopTimeMs"`
	DurationMs  int       `json:"durationMs"`
	Words       []VLFWord `json:"words"`
}

type VLFLattice map[string]VLFUtterance

/** read from a file the VLF Lattice */
func ParseVLF(filename string) (VLFLattice, error) {
	lattice := make(map[string]VLFUtterance)
	fileContents, err := ioutil.ReadFile(filename)
	err = json.Unmarshal(fileContents, lattice)
	if err != nil {
		return nil, err
	}
	return lattice, err
}
