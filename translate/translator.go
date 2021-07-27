package translate

type Translator struct {
	strict bool
}

func NewTranslator(strict bool) Translator {
	return Translator{strict: strict}
}

func StrictTranslator() Translator {
	return NewTranslator(true)
}
