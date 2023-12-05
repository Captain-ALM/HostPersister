package hosts

import "strings"

func NewHostsEntry(lineIn string) Entry {
	trLineIn := strings.ReplaceAll(strings.Trim(lineIn, "\r\n"), "	", " ")
	lineSplt := strings.Split(trLineIn, " ")
	if len(lineSplt) == 1 && strings.HasPrefix(trLineIn, "#") {
		return Entry{
			IPAddress: "",
			Domains:   nil,
			comment:   trLineIn,
		}
	} else if len(lineSplt) > 1 {
		var theDomains []string
		for i := 1; i < len(lineSplt); i++ {
			if lineSplt[i] == "" {
				continue
			}
			if strings.HasPrefix(lineSplt[i], "#") {
				break
			}
			theDomains = append(theDomains, lineSplt[i])
		}
		theComment := ""
		theCommentStart := strings.Index(trLineIn, "#")
		if theCommentStart > -1 {
			theComment = trLineIn[theCommentStart:]
		}
		return Entry{
			IPAddress: lineSplt[0],
			Domains:   theDomains,
			comment:   theComment,
		}
	} else {
		return Entry{
			IPAddress: "",
			Domains:   nil,
			comment:   "",
		}
	}
}

type Entry struct {
	IPAddress string
	Domains   []string
	comment   string
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
	if e.IsFilled() {
		toReturn := []string{e.IPAddress}
		toReturn = append(toReturn, e.Domains...)
		if e.comment != "" {
			toReturn = append(toReturn, e.comment)
		}
		return strings.Join(toReturn, " ")
	}
	return e.comment
}
