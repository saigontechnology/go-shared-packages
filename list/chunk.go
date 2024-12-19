package list

func Chunk[E any](values []E, size int) [][]E {
	if len(values) == 0 || size <= 0 {
		return [][]E{}
	}
	var chunks [][]E
	for remaining := len(values); remaining > 0; remaining = len(values) {
		if remaining < size {
			size = remaining
		}
		// Only append chunk with cap of size (memory optimization aspect)
		chunks = append(chunks, values[:size:size])
		values = values[size:]
	}
	return chunks
}
