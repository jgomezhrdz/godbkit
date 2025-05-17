package criteriamanager

import (
	kit "godbkit/internal/kit"
	"strconv"
	"strings"
)

func CriteriaFromRequest(queryParam map[string][]string) Criteria {
	limit := getLimits(queryParam, "limit")
	offset := getLimits(queryParam, "offset")
	order := getOrder(queryParam)
	orderBy := getOrderBy(queryParam)
	filtros := parseToFilterArrays(queryParam)
	return NewCriteria(filtros, limit, offset, order, orderBy)
}

func getOrder(params map[string][]string) *string {
	var order *string
	orderElems := params["order"]

	if orderElems != nil && len(orderElems) > 0 {
		orderAux := (orderElems[0])
		order = &orderAux
	}

	return order
}

func getOrderBy(params map[string][]string) *string {
	var order *string
	orderElems := params["orderBy"]

	if orderElems != nil && len(orderElems) > 0 {
		orderAux := kit.CamelToUnderscore(orderElems[0])
		order = &orderAux
	}

	return order
}

func getLimits(params map[string][]string, field string) *int {
	limitValues, ok := params[field]
	if !ok || len(limitValues) == 0 {
		return nil
	}
	limit, err := strconv.Atoi(limitValues[0])
	if err != nil {
		return nil
	}

	return &limit
}

func parseToFilterArrays(queryParam map[string][]string) [][]Filter {
	var conditionArrays [][]Filter

	// Iterate over query parameters
	for key, values := range queryParam {
		if strings.Contains(key, "filter") {
			// Split key into parts to extract array indices
			parts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(key, "filters["), "]"), "][")
			// Extract array indices and create Filter instances
			if len(parts) == 3 {
				index1, _ := strconv.Atoi(parts[0])
				index2, _ := strconv.Atoi(parts[1])
				filterKey := parts[2]
				for _, value := range values {
					// Ensure the array is large enough
					for len(conditionArrays) <= index1 {
						conditionArrays = append(conditionArrays, nil)
					}
					for len(conditionArrays[index1]) <= index2 {
						conditionArrays[index1] = append(conditionArrays[index1], Filter{})
					}
					// Set Filter values
					switch filterKey {
					case "campo":
						parts := strings.Split(value, "->")
						field := kit.CamelToUnderscore(parts[0])
						if len(parts) > 1 {
							field += "->" + strings.Join((parts[1:]), "->")
						}
						conditionArrays[index1][index2].Campo = field

					case "operador":
						conditionArrays[index1][index2].Operador = value
					case "valor":
						conditionArrays[index1][index2].Valor = value
					}
				}
			}
		}
	}

	var result [][]Filter
	for _, condition := range conditionArrays {
		var filteredCondition []Filter
		for _, item := range condition {
			if item.Valor != "undefined" {
				filteredCondition = append(filteredCondition, item)
			}
		}
		if len(filteredCondition) > 0 {
			result = append(result, filteredCondition)
		}
	}

	return result
}
