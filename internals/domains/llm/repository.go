package llm

type Repository interface {
	Prompt(prompt string) (string, error)
	Embed(prompt string) ([]float32, error)
}
