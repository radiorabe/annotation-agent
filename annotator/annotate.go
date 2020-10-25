package annotator

// Annotator base interface
type Annotator interface {
	CreateAnnotations(url string) ([]string, error)
}

// DefaultAnnotator ...
type DefaultAnnotator struct{}

// CreateAnnotations  ...
func (a *DefaultAnnotator) CreateAnnotations(url string) ([]string, error) {
	return []string{}, nil
}
