package types

type Endianness string

const (
	LITTLE_ENDIAN Endianness = "le"
	BIG_ENDIAN    Endianness = "be"
)

const (
	UTF_8  string = "utf-8"
	UTF_16 string = "utf-16"
	UTF_32 string = "utf-32"
)
