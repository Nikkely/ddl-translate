package translate

type Translater interface {
	Run(text string) (string, error)
}
