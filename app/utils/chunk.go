package utils

func ChunkSlice[T any](items []T, size int) [][]T {
	var chunks [][]T
	// Iterar
	for i := 0; i < len(items); i += size {
		end := i + size
		// Verificar que el tamaño sea mayor al chunk
		if end > len(items) {
			end = len(items)
		}
		// Agregar el chunk
		chunks = append(chunks, items[i:end])
	}
	// Regresar chunks
	return chunks
}
