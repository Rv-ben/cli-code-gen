using Qdrant.Client.Grpc;

namespace AiCodeEditor.Cli.Services
{
    public class CodeSearchService
    {
        private readonly OllamaEmbeddingService _embeddingService;
        private readonly QdrantService _qdrantService;

        public record SearchResult(
            string FilePath,
            int StartLine,
            int EndLine,
            string Language,
            string Content,
            float Score
        );

        public CodeSearchService(OllamaEmbeddingService? ollamaEmbeddingService, QdrantService? qdrantService)
        {
            if (ollamaEmbeddingService == null)
            {
                throw new ArgumentNullException(nameof(ollamaEmbeddingService));
            }

            if (qdrantService == null)
            {
                throw new ArgumentNullException(nameof(qdrantService));
            }

            _embeddingService = ollamaEmbeddingService;
            _qdrantService = qdrantService;
        }

        public async Task<List<SearchResult>> SearchAsync(
            string query,
            int limit = 5,
            float scoreThreshold = 0.7f)
        {
            // Generate embedding for the search query
            var queryEmbedding = await _embeddingService.GetEmbeddingAsync(query);
            
            // Search in Qdrant
            var searchResults = await _qdrantService.SearchAsync(
                queryEmbedding,
                limit,
                scoreThreshold
            );

            // Convert results to our format
            return searchResults
                .Select(result => new SearchResult(
                    FilePath: result.Payload["file_path"].StringValue,
                    StartLine: int.Parse(result.Payload["start_line"].StringValue),
                    EndLine: int.Parse(result.Payload["end_line"].StringValue),
                    Language: result.Payload["language"].StringValue,
                    Content: result.Payload["content"].StringValue,
                    Score: result.Score
                ))
                .OrderByDescending(r => r.Score)
                .ToList();
        }

        public async Task<List<string>> SearchFilePathsAsync(
            string query,
            int limit = 5,
            float scoreThreshold = 0.7f,
            List<string>? excludedFilePaths = null)
        {
            var results = await SearchAsync(query, limit, scoreThreshold);

            // Remove excluded file paths
            if (excludedFilePaths != null)
            {
                results = results.Where(result => !excludedFilePaths.Contains(result.FilePath)).ToList();
            }

            return results
                .Select(r => r.FilePath)
                .Distinct()
                .ToList();
        }

        public async Task<Dictionary<string, List<SearchResult>>> SearchGroupedByFileAsync(
            string query,
            int limit = 5,
            float scoreThreshold = 0.7f,
            List<string>? excludedFilePaths = null)
        {
            var results = await SearchAsync(query, limit, scoreThreshold);

            // Remove excluded file paths
            if (excludedFilePaths != null)
            {
                results = results.Where(result => !excludedFilePaths.Contains(result.FilePath)).ToList();
            }

            return results
                .GroupBy(r => r.FilePath)
                .ToDictionary(
                    g => g.Key,
                    g => g.OrderBy(r => r.StartLine).ToList()
                );
        }
    }
} 