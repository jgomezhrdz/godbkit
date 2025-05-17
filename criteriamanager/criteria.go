package criteriamanager

type Criteria struct {
	filtros [][]Filter
	limit   *int
	offset  *int
	order   *string
	ordeby  *string
}

func NewCriteria(filtros [][]Filter, limit *int, offset *int, order, orderBy *string) Criteria {
	return Criteria{filtros: filtros, limit: limit, offset: offset, order: order, ordeby: orderBy}
}

func EmptyCriteria() Criteria {
	return Criteria{filtros: [][]Filter{}, limit: nil, offset: nil}
}

func (c Criteria) ADDFILTROS(filtros []Filter) Criteria {
	c.filtros = append(c.filtros, filtros)
	return c
}

func (c Criteria) GETFILTROS() [][]Filter {
	return c.filtros
}

func (c Criteria) GETLIMIT() *int {
	return c.limit
}

func (c Criteria) GETOFFSET() *int {
	return c.offset
}

func (c Criteria) GETORDER() string {
	var order string = ""
	if c.ordeby != nil {
		order = *c.ordeby
	}
	if c.order != nil {
		order += " " + *c.order
	}
	return order
}
