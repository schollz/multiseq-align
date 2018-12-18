package multiseqalign

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	align, err := New("testing/template.seq")
	assert.Nil(t, err)
	fmt.Println(len(align.Template.SequenceMap))

	err = align.AddSequence("testing/S13-SH6.seq")
	assert.Nil(t, err)

	err = align.AddSequence("testing/S13-SH7.seq")
	assert.Nil(t, err)

	err = align.AddSequence("testing/S13-SH8.seq")
	assert.Nil(t, err)

	err = align.AddSequence("testing/S13-SH10.seq")
	assert.Nil(t, err)

	err = align.AddSequence("testing/S13-SH11.seq")
	assert.Nil(t, err)

	err = align.AddSequence("testing/S13-SH12.seq")
	assert.Nil(t, err)

	err = align.AddSequence("testing/S13-SH13.seq")
	assert.Nil(t, err)

	err = align.AddSequence("testing/S13-SH14.seq")
	assert.Nil(t, err)

	err = align.AddSequence("testing/S13-SH152.seq")
	assert.Nil(t, err)

	err = align.AddSequence("testing/S13-SH162.seq")
	assert.Nil(t, err)

	align.Print()
}
