package u

import (
	"fmt"

	"slices"
)

func ToStr(data any) string {
	return fmt.Sprintf("%v", data)
}

func AllStrInList(list []string, check []string) bool {
	for i := range check {
		if !slices.Contains(list, check[i]) {
			return false
		}
	}
	return true
}
