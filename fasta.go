package multiseqalign

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Fasta struct {
	Header   string
	Sequence string
}

func ImportFasta(fname string) (f Fasta, err error) {
	file, err := os.Open(fname)
	if err != nil {
		return
	}
	defer file.Close()
	f = Fasta{}
	foundHeader := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.ToUpper(strings.TrimSpace(scanner.Text()))
		if string(line[0]) == ">" && len(line) > 1 {
			f.Header = line[1:]
			foundHeader = true
			continue
		}
		f.Sequence += line
	}
	err = scanner.Err()
	if err != nil {
		return
	}
	if !foundHeader {
		err = fmt.Errorf("could not find valid FASTA header")
	}
	return
}
