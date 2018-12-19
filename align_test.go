package multiseqalign

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	align, err := New("testing/template.seq")
	assert.Nil(t, err)
	fmt.Println(len(align.Template.SequenceMap))

	var wg sync.WaitGroup

	files := []string{"testing/S13-SH6.seq", "testing/S13-SH7.seq", "testing/S13-SH8.seq", "testing/S13-SH93.seq", "testing/S13-SH10.seq", "testing/S13-SH11.seq", "testing/S13-SH12.seq", "testing/S13-SH13.seq", "testing/S13-SH14.seq", "testing/S13-SH152.seq", "testing/S13-SH162.seq"}

	for _, fname := range files {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()
			seq, err := ImportSequence(fname)
			assert.Nil(t, err)
			seq, err = AlignSequences(align.Template, seq)
			assert.Nil(t, err)
			align.AddAlignedSequence(seq)
		}(fname)
	}

	wg.Wait()
	// err = align.AddSequence("testing/S13-SH6.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH7.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH8.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH93.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH10.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH11.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH12.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH13.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH14.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH152.seq")
	// assert.Nil(t, err)

	// err = align.AddSequence("testing/S13-SH162.seq")
	// assert.Nil(t, err)

	align.Print()
}
