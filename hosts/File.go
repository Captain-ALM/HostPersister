package hosts

import (
	"io"
	"os"
	"strings"
)

const readBufferSize = 8192

func NewHostsFile(filePath string) (*File, error) {
	theHostFile := &File{
		filePath: filePath,
	}
	err := theHostFile.ReadHostsFile()
	if err == nil {
		return theHostFile, nil
	}
	return nil, err
}

type File struct {
	filePath   string
	Entries    []Entry
	lineEnding string
}

func (f *File) ReadHostsFile() error {
	f.Entries = nil
	theFile, err := os.Open(f.filePath)
	if err != nil {
		return err
	}
	defer theFile.Close()
	var lenIn int
	f.lineEnding = ""
	theCBuffer := ""
	theBuffer := make([]byte, readBufferSize)
	for err == nil {
		lenIn, err = theFile.Read(theBuffer)
		if lenIn > 0 {
			theCBuffer += string(theBuffer[:lenIn])
			if f.lineEnding == "" {
				if strings.Contains(theCBuffer, "\r\n") {
					f.lineEnding = "\r\n"
				} else if strings.Contains(theCBuffer, "\r") {
					f.lineEnding = "\r"
				} else if strings.Contains(theCBuffer, "\n") {
					f.lineEnding = "\n"
				}
			}
			if f.lineEnding == "\r\n" {
				strings.ReplaceAll(theCBuffer, "\r\n", "\n")
			} else if f.lineEnding == "\r" {
				strings.ReplaceAll(theCBuffer, "\r", "\n")
			}
			splt := strings.Split(theCBuffer, "\n")
			for i := 0; i < len(splt)-1; i++ {
				f.Entries = append(f.Entries, NewHostsEntry(splt[i]))
			}
			theCBuffer = splt[len(splt)-1]
		}
	}
	if err != io.EOF {
		return err
	}
	if theCBuffer != "" {
		f.Entries = append(f.Entries, NewHostsEntry(theCBuffer))
	}
	return nil
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
