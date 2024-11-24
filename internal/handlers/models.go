package handlers

type Handler struct {
	Name       string `yaml:"name"`
	Persistent bool   `yaml:"persistent"`
}
