package data

type Option int

const (
	OPTION_COPY_FULL_PATH Option = iota
)

var OptionsMap = map[Option]string{
	OPTION_COPY_FULL_PATH: "COPY_FULL_PATH",
}
