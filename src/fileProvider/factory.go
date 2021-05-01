package fileProvider

func Factory(ext string) FileProvider {
	if ext == ".yaml" {
		return NewYAMLProvider()
	}
}
