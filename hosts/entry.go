package hosts

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var ErrInvalidEntry = errors.New("invalid entry")

type Entry struct {
	IPAddress string
	Domains   []string
	Comment   string
}

func (e Entry) IsFilled() bool {
	return e.IPAddress != "" && len(e.Domains) > 0
}

func (e Entry) HasDomain(domain string) bool {
	if !e.IsFilled() {
		return false
	}
	for _, c := range e.Domains {
		if strings.EqualFold(c, domain) {
			return true
		}
	}
	return false
}

func (e Entry) ToLine() string {
	var b strings.Builder
	filled := e.IsFilled()
	if filled {
		b.WriteString(e.IPAddress)
		for _, i := range e.Domains {
			b.WriteByte(' ')
			b.WriteString(i)
		}
	}
	if e.Comment != "" {
		b.Grow(len(e.Comment) + 2)
		if filled {
			b.WriteByte(' ')
		}
		b.WriteByte('#')
		b.WriteString(e.Comment)
	}
	return b.String()
}

func ParseEntryString(line string) (Entry, error) {
	return ParseEntry(strings.NewReader(line))
}

func ParseEntry(line io.Reader) (entry Entry, err error) {
	cr := NewCommentReader(line, '#')
	sc := bufio.NewScanner(cr)
	sc.Split(bufio.ScanWords)

	isFirst := true
	for sc.Scan() {
		t := sc.Text()
		if isFirst {
			entry.IPAddress = t
			isFirst = false
		} else {
			entry.Domains = append(entry.Domains, t)
		}
	}
	err = sc.Err()
	if err != nil {
		return
	}

	// invalid if the ip address is set but no domains are added
	if entry.IPAddress != "" && len(entry.Domains) < 1 {
		err = ErrInvalidEntry
		return
	}
	var cAll []byte
	cAll, err = io.ReadAll(cr.Comment())
	if len(cAll) >= 1 {
		entry.Comment = string(cAll[1:])
	}
	return
}
