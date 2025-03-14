using AiCodeEditor.Cli.Models;
using static AiCodeEditor.Cli.Services.CodebasePathChunkingService;

namespace AiCodeEditor.Cli.Services
{
    public class CodebaseIndexingService
    {
        private readonly CodebaseChunkingService _chunkingService;
        private readonly CodebasePathChunkingService _pathChunkingService;
        private readonly OllamaEmbeddingService _embeddingService;
        private readonly QdrantService _qdrantService;
        private readonly AppConfig _config;

        public delegate void ProgressCallback(string message, bool isNewLine = true);

        public CodebaseIndexingService(
            CodebaseChunkingService chunkingService,
            CodebasePathChunkingService pathChunkingService,
            OllamaEmbeddingService embeddingService,
            QdrantService qdrantService,
            AppConfig config)
        {
            _chunkingService = chunkingService;
            _pathChunkingService = pathChunkingService;
            _embeddingService = embeddingService;
            _qdrantService = qdrantService;
            _config = config;
        }

        public async Task IndexCodebaseAsync(string directory, ProgressCallback onProgress)
        {

            // Initialize collection
            onProgress("Initializing collection...");
            await _qdrantService.InitializeCollectionAsync();

            try
            {
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

                    await _qdrantService.UpsertCodeBaseChunkVectorAsync((ulong)i, embedding, payload);
                    
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

            try
            {
                // Chunk file paths
                onProgress("Chunking file paths...");
                var pathChunks = _pathChunkingService.ChunkCodebasePaths(directory);
                onProgress($"Found {pathChunks.Count} path chunks across different granularities");

                // Store file paths with embeddings
                onProgress("Generating embeddings and storing path chunks...");
                for (int i = 0; i < pathChunks.Count; i++)
                {
                    var chunk = pathChunks[i];
                    
                    // Create embedding context based on chunk type
                    var embeddingContext = chunk.Type switch
                    {
                        ChunkType.FullPath => $"Full path: {chunk.Segment}",
                        ChunkType.DirectoryOnly => $"Directory structure: {chunk.Segment}",
                        ChunkType.FileNameOnly => $"File name: {chunk.Segment}",
                        ChunkType.DirectorySegment => $"Directory segment at depth {chunk.SegmentDepth}: {chunk.Segment}",
                        _ => chunk.Segment
                    };

                    var embedding = await _embeddingService.GetEmbeddingAsync(embeddingContext);
                    var payload = CodebasePathChunkingService.ToFieldMap(chunk);
                    var payloadDict = payload.ToDictionary(
                        kvp => kvp.Key,
                        kvp => kvp.Value.ToString()
                    );

                    await _qdrantService.UpsertFilePathVectorAsync((ulong)(i + 1), embedding, payloadDict);

                    if ((i + 1) % 10 == 0) // Progress indicator
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