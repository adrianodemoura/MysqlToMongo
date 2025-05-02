package migration

import (
	"runtime"
)

// Função para calcular o limite de memória
func calculateMemoryLimit() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	totalMemory := m.Sys
	memoryLimit := int64(float64(totalMemory) * 0.7)
	return memoryLimit
}

// Função para dividir o trabalho entre workers
func SplitWork(totalRecords int64, numWorkers int) []struct{ Start, End int64 } {
	chunkSize := totalRecords / int64(numWorkers)
	chunks := make([]struct{ Start, End int64 }, numWorkers)

	for i := 0; i < numWorkers; i++ {
		chunks[i].Start = int64(i)*chunkSize + 1 // Começa do 1 para não pular registros
		if i == numWorkers-1 {
			chunks[i].End = totalRecords
		} else {
			chunks[i].End = chunks[i].Start + chunkSize - 1
		}
	}

	return chunks
}
