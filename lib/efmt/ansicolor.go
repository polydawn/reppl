package efmt

import (
	"bytes"
)

func Ansi(codes ...[]byte) []byte {
	return bytes.Join([][]byte{
		csi,
		bytes.Join(codes, []byte{';'}),
		cmd_sgr,
	}, []byte{})
}

func AnsiWrap(str string, codes ...[]byte) []byte {
	return bytes.Join([][]byte{
		Ansi(codes...),
		[]byte(str),
		Ansi(Ansi_reset),
	}, []byte{})
}

var (
	csi                   = []byte("\033[")
	cmd_sgr               = []byte("m") // "Select Graphic Rendition"
	Ansi_textBlack        = []byte("30")
	Ansi_textDarkGray     = []byte("1;30")
	Ansi_textRed          = []byte("31")
	Ansi_textBrightRed    = []byte("1;31")
	Ansi_textGreen        = []byte("32")
	Ansi_textBrightGreen  = []byte("1;32")
	Ansi_textYellow       = []byte("33")
	Ansi_textBrightYellow = []byte("1;33")
	Ansi_textBlue         = []byte("34")
	Ansi_textBrightBlue   = []byte("1;34")
	Ansi_textPurple       = []byte("35")
	Ansi_textBrightPurple = []byte("1;35")
	Ansi_textCyan         = []byte("36")
	Ansi_textBrightCyan   = []byte("1;36")
	Ansi_textGray         = []byte("37")
	Ansi_textWhite        = []byte("1;37")
	Ansi_backgroundBlack  = []byte("40")
	Ansi_backgroundRed    = []byte("41")
	Ansi_backgroundGreen  = []byte("42")
	Ansi_backgroundYellow = []byte("43")
	Ansi_backgroundBlue   = []byte("44")
	Ansi_backgroundPurple = []byte("45")
	Ansi_backgroundCyan   = []byte("46")
	Ansi_backgroundGray   = []byte("47")
	Ansi_underline        = []byte("4")
	Ansi_blink            = []byte("5")
	Ansi_reset            = []byte{}
)
