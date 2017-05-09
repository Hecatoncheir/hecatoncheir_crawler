package crawler

type Company struct {
	ID         string
	Iri        string
	Name       string
	Categories []string
}

type EntityConfig struct {
	Company
}
