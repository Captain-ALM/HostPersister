package hosts

import (
	"io"
	"os"
	"strings"
)

const readBufferSize = 8192

func NewHostsFile(filePath string) (*File, error) {
	theFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer theFile.Close()
	var theEntries []Entry
	var lenIn int
	lineEnding := ""
	theCBuffer := ""
	theBuffer := make([]byte, readBufferSize)
	for err == nil {
		lenIn, err = theFile.Read(theBuffer)
		if lenIn > 0 {
			theCBuffer += string(theBuffer[:lenIn])
			if lineEnding == "" {
				if strings.Contains(theCBuffer, "\r\n") {
					lineEnding = "\r\n"
				} else if strings.Contains(theCBuffer, "\r") {
					lineEnding = "\r"
				} else if strings.Contains(theCBuffer, "\n") {
					lineEnding = "\n"
				}
			}
			if lineEnding == "\r\n" {
				strings.ReplaceAll(theCBuffer, "\r\n", "\n")
			} else if lineEnding == "\r" {
				strings.ReplaceAll(theCBuffer, "\r", "\n")
			}
			splt := strings.Split(theCBuffer, "\n")
			for i := 0; i < len(splt)-1; i++ {
				theEntries = append(theEntries, NewHostsEntry(splt[i]))
			}
			theCBuffer = splt[len(splt)-1]
		}
	}
	if err != io.EOF {
		return nil, err
	}
	if theCBuffer != "" {
		theEntries = append(theEntries, NewHostsEntry(theCBuffer))
	}
	return &File{
		filePath:   filePath,
		Entries:    theEntries,
		lineEnding: lineEnding,
	}, nil
}

type File struct {
	filePath   string
	Entries    []Entry
	lineEnding string
}

func (f File) WriteHostsFile() error {
	theFile, err := os.OpenFile(f.filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer theFile.Close()
	for _, entry := range f.Entries {
		_, err = theFile.WriteString(entry.ToLine() + f.lineEnding)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f File) HasDomain(domain string) bool {
	for _, entry := range f.Entries {
		if entry.HasDomain(domain) {
			return true
		}
	}
	return false
}

func (f File) indexDomainSingleton(domain string) int {
	for i, entry := range f.Entries {
		if len(entry.Domains) == 1 && entry.HasDomain(domain) {
			return i
		}
	}
	return -1
}

func (f File) HasDomainSingleton(domain string) bool {
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
