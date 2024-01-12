package hosts

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"unicode"
)

func TestParseEntry(t *testing.T) {
	entry, err := ParseEntry(strings.NewReader("127.0.0.1 myself.local another.local # this is a test comment"))
	assert.NoError(t, err)
	assert.Equal(t, Entry{
		IPAddress: "127.0.0.1",
		Domains:   []string{"myself.local", "another.local"},
		Comment:   " this is a test comment",
	}, entry)

	entry, err = ParseEntry(strings.NewReader("127.0.0.1 test.local #another comment"))
	assert.NoError(t, err)
	assert.Equal(t, Entry{
		IPAddress: "127.0.0.1",
		Domains:   []string{"test.local"},
		Comment:   "another comment",
	}, entry)

	entry, err = ParseEntry(strings.NewReader("127.0.0.1     "))
	assert.EqualError(t, err, "invalid entry")
}

func TestEntry_IsFilled(t *testing.T) {
	assert.True(t, Entry{IPAddress: "127.0.0.1", Domains: []string{"myself.local"}}.IsFilled())
	assert.False(t, Entry{}.IsFilled())
	assert.False(t, Entry{IPAddress: "127.0.0.1"}.IsFilled())
}

func TestEntry_HasDomain(t *testing.T) {
	assert.True(t, Entry{IPAddress: "127.0.0.1", Domains: []string{"myself.local", "another.local"}}.HasDomain("myself.local"))
	assert.True(t, Entry{IPAddress: "127.0.0.1", Domains: []string{"myself.local", "another.local"}}.HasDomain("MYSELF.local"))
	assert.False(t, Entry{IPAddress: "127.0.0.1", Domains: []string{"notme.local", "another.local"}}.HasDomain("myself.local"))
}

func TestEntry_ToLine(t *testing.T) {
	assert.Equal(t, "127.0.0.1 myself.local another.local # this is a test comment", Entry{IPAddress: "127.0.0.1", Domains: []string{"myself.local", "another.local"}, Comment: " this is a test comment"}.ToLine())
	assert.Equal(t, "# this is a test comment", Entry{IPAddress: "", Domains: nil, Comment: " this is a test comment"}.ToLine())
	assert.Equal(t, "127.0.0.1 myself.local another.local", Entry{IPAddress: "127.0.0.1", Domains: []string{"myself.local", "another.local"}}.ToLine())
}

func FuzzParseEntry(f *testing.F) {
	f.Add("127.0.0.1", "myself.local", "another.local", "this is a test comment")
	f.Fuzz(func(t *testing.T, a, b, c, d string) {
		for _, i := range []string{a, b, c} {
			i = strings.TrimSpace(i)
			if i == "" || strings.ContainsFunc(i, func(r rune) bool {
				return r == '#' || unicode.IsSpace(r)
			}) {
				t.Skip()
			}
		}

		entry := Entry{IPAddress: a, Domains: []string{b, c}, Comment: d}
		e2, err := ParseEntryString(entry.ToLine())
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, Entry{
			IPAddress: strings.TrimSpace(a),
			Domains:   []string{strings.TrimSpace(b), strings.TrimSpace(c)},
			Comment:   d,
		}, e2)
	})
}
