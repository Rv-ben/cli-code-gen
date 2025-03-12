package services

import (
	"strings"
)

// CodeChunkingService handles splitting code into manageable chunks
type CodeChunkingService struct {
	defaultChunkSize int
}

// NewCodeChunkingService creates a new instance of the code chunking service
func NewCodeChunkingService(defaultChunkSize int) *CodeChunkingService {
	if defaultChunkSize <= 0 {
		defaultChunkSize = 1000 // Default chunk size if invalid value provided
	}

	return &CodeChunkingService{
		defaultChunkSize: defaultChunkSize,
	}
}

// SplitCodeIntoChunks splits code into chunks of approximately the given size
func (s *CodeChunkingService) SplitCodeIntoChunks(code string, chunkSize int) []string {
	if chunkSize <= 0 {
		chunkSize = s.defaultChunkSize
	}

	lines := strings.Split(code, "\n")
	chunks := make([]string, 0)

	currentChunk := ""
	for _, line := range lines {
		// If adding this line would exceed chunk size, start a new chunk
		if len(currentChunk)+len(line)+1 > chunkSize && len(currentChunk) > 0 {
			chunks = append(chunks, currentChunk)
			currentChunk = line
		} else {
			if len(currentChunk) > 0 {
				currentChunk += "\n"
			}
			currentChunk += line
		}
	}

	// Add the last chunk if it's not empty
	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

// ChunkCodeWithMetadata splits code and returns chunks with updated metadata
func (s *CodeChunkingService) ChunkCodeWithMetadata(
	code string,
	baseMetadata map[string]interface{},
	chunkSize int,
) ([]string, []map[string]interface{}) {
	chunks := s.SplitCodeIntoChunks(code, chunkSize)
	metadataList := make([]map[string]interface{}, len(chunks))

	for i := range chunks {
		// Copy base metadata for each chunk
		chunkMetadata := copyMetadata(baseMetadata)
		chunkMetadata["chunk_index"] = i
		chunkMetadata["total_chunks"] = len(chunks)
		metadataList[i] = chunkMetadata
	}

	return chunks, metadataList
}

// Helper function to copy metadata map
func copyMetadata(metadata map[string]interface{}) map[string]interface{} {
	if metadata == nil {
		return make(map[string]interface{})
	}

	copy := make(map[string]interface{})
	for k, v := range metadata {
		copy[k] = v
	}

	return copy
}
