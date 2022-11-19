package parse

func (m *Model) Migration(deleteColumn bool) Migrationer {
	return m.Storage.Migration(m, deleteColumn)
}
