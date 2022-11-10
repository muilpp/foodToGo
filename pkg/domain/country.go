package domain

type Country struct {
	name string
}

func NewCountry(name string) *Country {
	return &Country{name}
}

func (c Country) GetName() string {
	return c.name
}
