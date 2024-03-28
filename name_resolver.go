package zinc

type NameResolver interface {
	ResolveTableName(structName string) string
	ResolveColumnName(structName, fieldName string) string
}

var DefaultNameResolver NameResolver = defaultNameResolver{}

type defaultNameResolver struct{}

func (defaultNameResolver) ResolveTableName(structName string) string {
	return structName
}

func (defaultNameResolver) ResolveColumnName(structName, fieldName string) string {
	return fieldName
}
