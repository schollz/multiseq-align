package multiseqalign

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/schollz/progressbar"
)

type Alignment struct {
	Template  Sequence
	Sequences []Sequence
	sync.RWMutex
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

func ImportSequence(sequenceFile string) (seq Sequence, err error) {
	fasta, err := ImportFasta(sequenceFile)
	if err != nil {
		return
	}
	seq = Sequence{Fasta: fasta, SequenceMap: make(map[int]rune)}
	for i, r := range seq.Fasta.Sequence {
		seq.SequenceMap[i] = r
	}
	return
}

func (a *Alignment) AddAlignedSequence(s Sequence) {
	a.Lock()
	defer a.Unlock()
	a.Sequences = append(a.Sequences, s)
}

// AlignSequences will align sequence2 to sequence1 and return the aligned sequence
func AlignSequences(s1, s2 Sequence) (sNew Sequence, err error) {
	fmt.Printf("\naligning %s\n", s2.Fasta.Header)

	// align sequence
	sequenceLength := len(s2.Fasta.Sequence)
	templateLength := len(s1.Fasta.Sequence)
	bestMatches := 0
	bestI := 0
	bar := progressbar.New(2 * (templateLength - sequenceLength))
	for i := 0; i < templateLength-sequenceLength; i++ {
		bar.Add(1)
		// count how many match
		matches := 0
		for j := range s2.SequenceMap {
			if s2.SequenceMap[j] == s1.SequenceMap[j+i] {
				matches++
			}
		}
		if matches > bestMatches {
			bestMatches = matches
			bestI = i
		}
	}

	revComp := ReverseComplement(s2)
	bestMatchesRevComp := 0
	bestIRevComp := 0
	for i := 0; i < templateLength-sequenceLength; i++ {
		bar.Add(1)
		// count how many match
		matches := 0
		for j := range revComp.SequenceMap {
			if revComp.SequenceMap[j] == s1.SequenceMap[j+i] {
				matches++
			}
		}
		if matches > bestMatchesRevComp {
			bestMatchesRevComp = matches
			bestIRevComp = i
		}
	}
	sNew = s2
	if bestMatchesRevComp > bestMatches {
		fmt.Println("using reverse complement")
		sNew = revComp
		sNew.Offset = bestIRevComp
	} else {
		sNew.Offset = bestI
	}

	return
}

func (a *Alignment) Print() {
	os.Remove("alignment.txt")
	f, err := os.Create("alignment.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

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
		for lineI, line := range lines {
			line = strings.Replace(line, " ", "", -1)
			if len(line) == 0 {
				continue
			}
			if lineI == 1 {
				f.WriteString(fmt.Sprintf("%5d ", i+1))
			} else if lineI == 0 {
				f.WriteString("     *")
			} else {
				f.WriteString("      ")
			}
			f.WriteString(line)
			if lineI == 1 {
				f.WriteString(fmt.Sprintf(" %5d", i+61))
			} else if lineI == 0 {
				f.WriteString("*     ")
			} else {
				f.WriteString("      ")
			}
			f.WriteString("\n")
		}
		f.WriteString("\n\n")
		i += 60
		if i > len(templateAligned) {
			break
		}
	}

}
