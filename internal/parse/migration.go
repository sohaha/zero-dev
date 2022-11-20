package parse

func (m *Modeler) Migration() Migrationer {
	return m.Storage.Migration(m)
}
