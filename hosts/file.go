package hosts

import (
	"bufio"
	"os"
)

func NewHostsFile(filePath string) (*File, error) {
	theHostFile := &File{filePath: filePath}
	err := theHostFile.ReadHostsFile()
	if err == nil {
		return theHostFile, nil
	}
	return nil, err
}

type File struct {
	filePath string
	Entries  []Entry
	LF       string
}

func (f *File) ReadHostsFile() error {
	f.Entries = nil
	theFile, err := os.Open(f.filePath)
	if err != nil {
		return err
	}
	defer theFile.Close()

	sc := bufio.NewScanner(theFile)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		t := sc.Text()
		entry, err := ParseEntryString(t)
		if err != nil {
			return err
		}
		f.Entries = append(f.Entries, entry)
	}
	return sc.Err()
}

func (f *File) WriteHostsFile() error {
	// default LF to \n
	if f.LF == "" {
		f.LF = "\n"
	}
	theFile, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer theFile.Close()
	for _, entry := range f.Entries {
		_, err = theFile.WriteString(entry.ToLine() + f.LF)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *File) HasDomain(domain string) bool {
	for _, entry := range f.Entries {
		if entry.HasDomain(domain) {
			return true
		}
	}
	return false
}

func (f *File) indexDomainSingleton(domain string) int {
	for i, entry := range f.Entries {
		if len(entry.Domains) == 1 && entry.HasDomain(domain) {
			return i
		}
	}
	return -1
}

func (f *File) HasDomainSingleton(domain string) bool {
	return f.indexDomainSingleton(domain) > -1
}

func (f *File) OverwriteDomainSingleton(domain string, ipAddress string) {
	idx := f.indexDomainSingleton(domain)
	if idx == -1 {
		f.Entries = append(f.Entries, Entry{
			IPAddress: ipAddress,
			Domains:   []string{domain},
		})
	} else {
		theEntry := f.Entries[idx]
		theEntry.IPAddress = ipAddress
		f.Entries = append(append(f.Entries[:idx], theEntry), f.Entries[idx+1:]...)
	}
}
