package zinc

type Options struct {
	// dialect
	Dialect Dialect

	// mapper
	Mapper Mapper

	// log
	Logger           Logger
	LogFormatter     LogFormatter
	LogBound         bool
	LogSuccess       bool
	LogSlowThreshold int64
}

type OptionsModifier func(*Options)

func modifyOptions(opts *Options, modifier OptionsModifier) *Options {
	opts1 := fromPtr(opts)
	if modifier != nil {
		modifier(&opts1)
	}
	return &opts1
}
