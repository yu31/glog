package glog

import (
	"unicode/utf8"

	"github.com/yu31/glog/pkg/buffer"
)

const hex = "0123456789abcdef"

var noEscapeTable = [256]bool{}

func init() {
	for i := 0; i <= 0x7e; i++ {
		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
	}
}

// AppendStringEscape encodes the input string to buffer.Buffer
//
// The operation loops though each byte in the string looking
// for characters that need json or utf8 encoding. If the string
// does not need encoding, then the string is appended in it's
// entirety to the byte slice.
// If we encounter a byte that does need encoding, switch up
// the operation and perform a byte-by-byte read-encode-append.
func AppendStringEscape(buf *buffer.Buffer, s string) {
	// Loop through each character in the string.
	for i := 0; i < len(s); i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if !noEscapeTable[s[i]] {
			// We encountered a character that needs to be encoded. Switch
			// to complex version of the algorithm.
			appendStringComplex(buf, s, i)
			return
		}
	}
	// The string has no need for encoding an therefore is directly
	// appended to the byte slice.
	buf.AppendString(s)
}

// appendStringComplex is used by appendString to take over an in
// progress JSON string encoding that encountered a character that needs
// to be encoded.
func appendStringComplex(buf *buffer.Buffer, s string, i int) {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				// In case of error, first append previous simple characters to
				// the byte slice if any and append a remplacement character code
				// in place of the invalid sequence.
				if start < i {
					buf.AppendString(s[start:i])
				}
				buf.AppendString(`\ufffd`)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if noEscapeTable[b] {
			i++
			continue
		}
		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			buf.AppendString(s[start:i])
		}
		switch b {
		case '"', '\\':
			buf.AppendByte('\\')
			buf.AppendByte(b)
		case '\b':
			buf.AppendByte('\\')
			buf.AppendByte('b')
		case '\f':
			buf.AppendByte('\\')
			buf.AppendByte('f')
		case '\n':
			buf.AppendByte('\\')
			buf.AppendByte('n')
		case '\r':
			buf.AppendByte('\\')
			buf.AppendByte('r')
		case '\t':
			buf.AppendByte('\\')
			buf.AppendByte('t')
		default:
			// Encode bytes < 0x20, except for the escape sequences above.
			buf.AppendString(`\u00`)
			buf.AppendByte(hex[b>>4])
			buf.AppendByte(hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		buf.AppendString(s[start:])
	}
}
