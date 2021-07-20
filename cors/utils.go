package cors

func resolveOption(option ...Option) Option {
	switch len(option) {
	case 0:
		return DefaultOption()
	default:
		return option[0]
	}
}
