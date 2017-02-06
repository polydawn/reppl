package efmt

import (
	"bufio"
	"bytes"
	"io"
)

var (
	non = []byte{}
	tab = []byte{'\t'}
	br  = []byte{'\n'}
)

// Proxies content by line, calling the proxied writer exactly once per line,
// applying no modifications to the content.
// In other words, this *just* batches write calls.
func LineBufferingWriter(w io.Writer) io.Writer {
	return LinePrefixingWriter(w, non)
}

// Proxies content by line, indenting each line with a tab character at start.
// Write calls are proxied once per line, and buffered until linebreak.
func LineIndentingWriter(w io.Writer) io.Writer {
	return LinePrefixingWriter(w, tab)
}

// Proxies content by line, prefixing each line with the given byte sequence.
// Write calls are proxied once per line, and buffered until linebreak.
func LinePrefixingWriter(w io.Writer, prefix []byte) io.Writer {
	return &Reframer{
		Delegate:  w,
		SplitFunc: bufio.ScanLines,
		Prefix:    prefix,
		Suffix:    br,
	}
}

// Proxies content by line, prefixing and suffixing each line with the given byte sequences.
// Write calls are proxied once per line, and buffered until linebreak.
func LineFlankingWriter(w io.Writer, prefix, suffix []byte) io.Writer {
	return &Reframer{
		Delegate:  w,
		SplitFunc: bufio.ScanLines,
		Prefix:    prefix,
		Suffix:    append(suffix[:], '\n'),
	}
}

var _ io.Writer = &Reframer{}

// Reframer implements `io.Writer`, and delegates to another Writer, buffering
// bytes into batches as defined by a `bufio.SplitFunc`, and applying
// a prefix and suffix byte slice before flushing each batch.
type Reframer struct {
	Delegate  io.Writer
	SplitFunc bufio.SplitFunc
	Prefix    []byte
	Suffix    []byte

	rem []byte
}

func (rfrm *Reframer) Write(b []byte) (int, error) {
	rfrm.rem = append(rfrm.rem, b...)
	for len(rfrm.rem) > 0 { // if loop until the buffer is exhausted, or another cond breaks out
		adv, tok, err := rfrm.SplitFunc(rfrm.rem, false)
		if err != nil {
			return len(b), err
		}
		if adv == 0 { // when we no longer have a full chunk, return
			return len(b), nil
		}
		// join all the things we're about to write, because we want them emitted as an atom.
		// (this may matter if the writer we're pushing into internally mutexes for sharing, for example.)
		rfrm.Delegate.Write(bytes.Join([][]byte{
			rfrm.Prefix,
			tok,
			rfrm.Suffix,
		}, []byte{}))
		rfrm.rem = rfrm.rem[adv:]
	}
	// we always state the entire range of bytes provided was written,
	//  because it is... if in buffer.  but we defintely don't need it re-sent.
	return len(b), nil
}
