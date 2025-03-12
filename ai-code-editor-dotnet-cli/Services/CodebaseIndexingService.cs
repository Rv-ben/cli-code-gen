using AiCodeEditor.Cli.Models;

namespace AiCodeEditor.Cli.Services
{
    public class CodebaseIndexingService
    {
        private readonly CodebaseChunkingService _chunkingService;
        private readonly OllamaEmbeddingService _embeddingService;
        private readonly QdrantService _qdrantService;
        private readonly AppConfig _config;

        public delegate void ProgressCallback(string message, bool isNewLine = true);

        public CodebaseIndexingService(
            CodebaseChunkingService chunkingService,
            OllamaEmbeddingService embeddingService,
            QdrantService qdrantService,
            AppConfig config)
        {
            _chunkingService = chunkingService;
            _embeddingService = embeddingService;
            _qdrantService = qdrantService;
            _config = config;
        }

        public async Task IndexCodebaseAsync(string directory, ProgressCallback onProgress)
        {
            try
            {
                // Initialize collection
                onProgress("Initializing collection...");
                await _qdrantService.InitializeCollectionAsync();

                // Chunk the codebase
                onProgress("Chunking codebase...");
                var chunks = await _chunkingService.ChunkCodebaseAsync(directory);
                onProgress($"Found {chunks.Count} chunks");

                // Store chunks with embeddings
                onProgress("Generating embeddings and storing chunks...");
                for (int i = 1; i < chunks.Count + 1; i++)
                {
                    var chunk = chunks[i - 1];
                    var embedding = await _embeddingService.GetEmbeddingAsync(chunk.Content);
                    var payload = new Dictionary<string, string>
                    {
                        { "file_path", chunk.FilePath },
                        { "start_line", chunk.StartLine.ToString() },
                        { "end_line", chunk.EndLine.ToString() },
                        { "language", chunk.Language },
                        { "content", chunk.Content }
                    };

                    await _qdrantService.UpsertVectorAsync((ulong)i, embedding, payload);
                    
                    if (i % 10 == 0) // Progress indicator
                    {
                        onProgress(".", false);
                    }
                }
                onProgress("\nStorage complete!");
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to index codebase: {ex.Message}", ex);
            }
        }
    }
} 