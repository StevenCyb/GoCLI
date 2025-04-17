package restriction

// IsCliOption is an interface that marks a type as a CLI option.
type IsCliOption interface {
	isCliOption()
}

// IsOptionOption is an interface that marks a type as an option.
type IsOptionOption interface {
	isOptionOption()
}

// IsCommandOption is an interface that marks a type as a command option.
type IsCommandOption interface {
	isCommandOption()
}

// IsArgumentOption is an interface that marks a type as an argument option.
type IsArgumentOption interface {
	isArgumentOption()
}
