package multiseqalign

import (
	"fmt"
	"strings"

	"github.com/schollz/progressbar"
)

type Alignment struct {
	Template  Sequence
	Sequences []Sequence
}

type Sequence struct {
	Fasta       Fasta
	SequenceMap map[int]rune
	Offset      int
}

func New(TemplateFile string) (a *Alignment, err error) {
	template, err := ImportFasta(TemplateFile)
	if err != nil {
		return
	}
	a = new(Alignment)
	a.Template = Sequence{Fasta: template}
	a.Template.SequenceMap = make(map[int]rune)
	for i, r := range a.Template.Fasta.Sequence {
		a.Template.SequenceMap[i] = r
	}
	a.Sequences = []Sequence{}
	return
}

func (a *Alignment) AddSequence(sequenceFile string) (err error) {
	fasta, err := ImportFasta(sequenceFile)
	if err != nil {
		return
	}
	a.Sequences = append(a.Sequences, Sequence{Fasta: fasta, SequenceMap: make(map[int]rune)})
	sI := len(a.Sequences) - 1
	for i, r := range a.Sequences[sI].Fasta.Sequence {
		a.Sequences[sI].SequenceMap[i] = r
	}

	// align sequence
	// TODO: determine if reverse complement is a better fit and automatically switch sequence
	sequenceLength := len(a.Sequences[sI].Fasta.Sequence)
	templateLength := len(a.Template.Fasta.Sequence)
	bestMatches := 0
	bestI := 0
	bar := progressbar.New(2 * (templateLength - sequenceLength))
	for i := 0; i < templateLength-sequenceLength; i++ {
		bar.Add(1)
		// count how many match
		matches := 0
		for j := range a.Sequences[sI].SequenceMap {
			if a.Sequences[sI].SequenceMap[j] == a.Template.SequenceMap[j+i] {
				matches++
			}
		}
		if matches > bestMatches {
			bestMatches = matches
			bestI = i
		}
	}

	revComp := ReverseComplement(a.Sequences[sI])
	bestMatchesRevComp := 0
	bestIRevComp := 0
	for i := 0; i < templateLength-sequenceLength; i++ {
		bar.Add(1)
		// count how many match
		matches := 0
		for j := range revComp.SequenceMap {
			if revComp.SequenceMap[j] == a.Template.SequenceMap[j+i] {
				matches++
			}
		}
		if matches > bestMatchesRevComp {
			bestMatchesRevComp = matches
			bestIRevComp = i
		}
	}
	if bestMatchesRevComp > bestMatches {
		a.Sequences[sI] = revComp
		a.Sequences[sI].Offset = bestIRevComp
	} else {
		a.Sequences[sI].Offset = bestI
	}

	return
}

func (a *Alignment) Print() {
	// create aligned template
	templateAligned := make(map[int][]rune)
	for i := range a.Template.SequenceMap {
		templateAligned[i] = []rune{a.Template.SequenceMap[i]}
	}

	for _, s := range a.Sequences {
		for i := range s.SequenceMap {
			templateAligned[i+s.Offset] = append(templateAligned[i+s.Offset], s.SequenceMap[i])
		}
	}

	i := 0
	for {
		maxLines := 0
		for j := 0; j < 60; j++ {
			if j+i > len(templateAligned) {
				break
			}
			if len(templateAligned[j+i]) > maxLines {
				maxLines = j
			}
		}

		lines := make([]string, maxLines+1)

		for j := 0; j < 60; j++ {
			lastI := 0
			hasMatch := false
			for runeI, r := range templateAligned[j+i] {
				if runeI > 0 {
					if string(r) == string(templateAligned[j+i][0]) {
						hasMatch = true
					}
				}
				lines[runeI+1] += string(r)
				lastI = runeI
			}
			if hasMatch {
				lines[0] += "*"
			} else {
				lines[0] += "-"
			}
			for k := lastI; k < len(lines); k++ {
				lines[k] += " "
			}
		}
		for _, line := range lines {
			line = strings.Replace(line, " ", "", -1)
			if len(line) == 0 {
				continue
			}
			fmt.Println(line)
		}
		fmt.Println(" ")
		i += 60
		if i > len(templateAligned) {
			break
		}
	}

}
