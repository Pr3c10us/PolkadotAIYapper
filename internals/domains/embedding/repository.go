package embedding

type Repository interface {
	AddEmbedding(embedding []float32, value string) error
	SimilarValuesExist(embedding []float32) (*bool, error)
}
