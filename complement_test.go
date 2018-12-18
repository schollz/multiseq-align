package multiseqalign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseComplement(t *testing.T) {
	seq := Sequence{Fasta: Fasta{Sequence: "NATCG", Header: "4BP"}}
	seq2 := ReverseComplement(seq)
	assert.Equal(t, "CGATN", seq2.Fasta.Sequence)
}
