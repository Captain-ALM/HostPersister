package hosts

import (
	"bytes"
	"io"
)

type CommentReader struct {
	r    io.Reader
	over []byte
	mark byte
	hit  bool
}

var _ io.Reader = &CommentReader{}

func NewCommentReader(r io.Reader, mark byte) *CommentReader {
	return &CommentReader{r, nil, mark, false}
}

func (c *CommentReader) Read(p []byte) (n int, err error) {
	if c.hit {
		return 0, io.EOF
	}
	n, err = c.r.Read(p)
	if err != nil {
		return
	}
	n2 := bytes.IndexByte(p[:n], c.mark)
	if n2 != -1 {
		c.over = p[n2:n]
		n = n2
		err = io.EOF
	}
	return
}

func (c *CommentReader) Comment() io.Reader {
	return io.MultiReader(bytes.NewReader(c.over), c.r)
}
