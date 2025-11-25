package types

type Endianness string

const (
	LITTLE_ENDIAN Endianness = "le"
	BIG_ENDIAN    Endianness = "be"
)

const (
	UTF_8    string = "utf-8"
	UTF_16   string = "utf-16"
	UTF_16LE string = "utf-16le"
	UTF_16BE string = "utf-16be"
	UTF_32   string = "utf-32"
	UTF_32LE string = "utf-32le"
	UTF_32BE string = "utf-32be"
)
