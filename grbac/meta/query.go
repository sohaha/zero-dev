package meta

// Query defines the data structure of the query parameters
type Query Resource

// GetArguments is used to convert the current argument to a string slice
func (query *Query) GetArguments() []string {
	return []string{
		query.Host,
		query.Path,
		query.Method,
	}
}
