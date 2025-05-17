package criteriamanager

import (
	"fmt"
	"strings"
)

// Filter represents the structure of the filter parameters
type Filter struct {
	Campo    string
	Operador string
	Valor    string
}

func NewFilter(campo string, operador string, valor string) Filter {
	return Filter{Campo: campo, Operador: operador, Valor: valor}
}

func ParseConditions(params [][]Filter) (string, []interface{}) {
	var conditions []string
	var queryValues []interface{}
	for _, elem := range params {
		var subConditions []string
		for _, filter := range elem {
			condition, valor := generateCondition(&filter)
			if valor != "NULL" && valor != "NOT NULL" {
				queryValues = append(queryValues, valor)
			}
			if condition != "" {
				subConditions = append(subConditions, condition)
			}
		}
		// Combine conditions within the same subarray with OR
		subCondition := strings.Join(subConditions, " OR ")
		if len(subConditions) > 1 {
			subCondition = "(" + subCondition + ")"
		}
		conditions = append(conditions, subCondition)
	}
	// Combine conditions between arrays with AND
	return strings.Join(conditions, " AND "), queryValues
}

// generateCondition generates a GORM condition string based on the filter parameters
func generateCondition(filter *Filter) (string, any) {
	query := ""
	var valor any
	if filter.Valor == "true" || filter.Valor == "false" {
		valor = filter.Valor == "true"
	} else {
		valor = filter.Valor
	}
	sprintfTemplate := "%s"

	//parse json values
	result := convertToJSONPath(filter.Campo)
	switch filter.Operador {
	case "=":
		if filter.Valor == "NULL" || filter.Valor == "NOT NULL" {
			query = fmt.Sprintf(sprintfTemplate+" IS "+filter.Valor, result)
		} else {
			query = fmt.Sprintf(sprintfTemplate+" = ?", result)
		}
	case ">":
		query = fmt.Sprintf(sprintfTemplate+" > ?", result)
	case "<>":
		query = fmt.Sprintf(sprintfTemplate+" <> ?", result)
	case ">=":
		query = fmt.Sprintf(sprintfTemplate+" > ?", result)
	case "<":
		query = fmt.Sprintf(sprintfTemplate+" < ?", result)
	case "<=":
		query = fmt.Sprintf(sprintfTemplate+" <= ?", result)
	case "IS", "is":
		query = fmt.Sprintf(sprintfTemplate+" IS "+filter.Valor, result)
	case "LIKE", "like":
		valor = "%" + valor.(string) + "%"
		query = fmt.Sprintf("%s LIKE ?", result)
	}

	return query, valor
}

func convertToJSONPath(input string) any {
	parts := strings.Split(input, "->")
	var expression string
	const template = "JSON_EXTRACT(%s, '$.%s')"
	if len(parts) > 1 {
		for i, part := range parts {
			if i == 1 {
				var field string = parts[0]
				expression = fmt.Sprintf(template, field, part)
			}
			if i > 1 {
				expression = fmt.Sprintf(template, expression, part)
			}
		}
	} else {
		expression = parts[0]
	}

	return expression
}
