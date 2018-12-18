package multiseqalign

func ReverseComplement(s Sequence) (s2 Sequence) {
	newSequence := ""
	for i := range s.Fasta.Sequence {
		switch seq := string(s.Fasta.Sequence[i]); seq {
		case "A":
			newSequence = "T" + newSequence
		case "C":
			newSequence = "G" + newSequence
		case "G":
			newSequence = "C" + newSequence
		case "T":
			newSequence = "A" + newSequence
		default:
			newSequence = seq + newSequence
		}
	}
	s2 = Sequence{Fasta: Fasta{Sequence: newSequence, Header: s.Fasta.Header + "REVERSE COMPLEMENT "}}
	s2.SequenceMap = make(map[int]rune)
	for i, r := range s2.Fasta.Sequence {
		s2.SequenceMap[i] = r
	}
	return
}
