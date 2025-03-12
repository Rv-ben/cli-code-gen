using AiCodeEditor.Cli.Services;
using AiCodeEditor.Cli.Models;
using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;

namespace AiCodeEditor.Cli.Commands
{
    [Command("search", Description = "Create a vector database from current codebase and search it")]
    public class SearchCodeCommand : ICommand
    {
        private readonly CodebaseChunkingService _chunkingService;
        private readonly OllamaEmbeddingService _embeddingService;
        private readonly QdrantService _qdrantService;
        private readonly CodeSearchService _searchService;
        private readonly AppConfig _config;
        [CommandOption("query", 'q', Description = "Search query")]
        public required string Query { get; init; }

        [CommandOption("limit", 'l', Description = "Maximum number of results to return")]
        public int Limit { get; init; } = 5;

        public SearchCodeCommand(
            CodebaseChunkingService chunkingService,
            OllamaEmbeddingService embeddingService,
            QdrantService qdrantService,
            CodeSearchService searchService,
            AppConfig config)
        {
            _chunkingService = chunkingService;
            _embeddingService = embeddingService;
            _qdrantService = qdrantService;
            _searchService = searchService;
            _config = config;
        }

        public async ValueTask ExecuteAsync(IConsole console)
        {
            try
            {
                // Initialize collection
                console.Output.WriteLine($"Initializing collection...");
                await _qdrantService.InitializeCollectionAsync();

                // Chunk the codebase
                console.Output.WriteLine("Chunking codebase...");
                var chunks = await _chunkingService.ChunkCodebaseAsync(Directory.GetCurrentDirectory());
                console.Output.WriteLine($"Found {chunks.Count} chunks");

                // Store chunks with embeddings
                console.Output.WriteLine("Generating embeddings and storing chunks...");
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
                        console.Output.Write(".");
                    }
                }
                console.Output.WriteLine("\nStorage complete!");

                // Search
                console.Output.WriteLine($"\nSearching for: {Query}");
                var results = await _searchService.SearchGroupedByFileAsync(Query, Limit, _config.SearchThreshold);

                // Display results
                if (!results.Any())
                {
                    console.Output.WriteLine("No results found.");
                    return;
                }

                foreach (var (filePath, matches) in results)
                {
                    console.Output.WriteLine($"\nFile: {filePath}");
                    foreach (var match in matches)
                    {
                        console.Output.WriteLine($"  Lines {match.StartLine}-{match.EndLine} (Score: {match.Score:F2})");
                        console.Output.WriteLine("  Preview:");
                        var preview = match.Content.Length > 200 
                            ? match.Content[..200] + "..."
                            : match.Content;
                        console.Output.WriteLine($"  {preview.Replace("\n", "\n  ")}");
                        console.Output.WriteLine();
                    }
                }
            }
            catch (Exception ex)
            {
                console.Error.WriteLine($"Error: {ex.Message}");
                throw;
            }
        }
    }
} 