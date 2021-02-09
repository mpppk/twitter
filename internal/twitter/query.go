package twitter

import "strings"

func BuildQuery(query string, ignoreKeywords []string, excludes []string, filters []string) string {
	queries := []string{query}
	for _, keyword := range ignoreKeywords {
		queries = append(queries, "-"+keyword)
	}
	for _, exclude := range excludes {
		queries = append(queries, "exclude:"+exclude)
	}
	for _, filter := range filters {
		queries = append(queries, "filter:"+filter)
	}
	return strings.Join(queries, " ")
}
