package loader

func (l *Loader) loadModules() {
	if l.err != nil {
		return
	}
	l.loadModeler("./app/modules")
	l.loadViews("./app/modules")
}
