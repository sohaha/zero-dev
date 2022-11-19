package parse

func (m *Model) Migration() Migrationer {
	return m.Storage.Migration(m)
}
