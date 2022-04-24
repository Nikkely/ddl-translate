package translate

type Deepl struct {
	apiKey string
}

func NewDeepl(key string) *Deepl {
	return &Deepl{
		apiKey: key,
	}
}

func (d Deepl) Run(text string) (string, error) {
	// request deepl
	return "", nil
}
