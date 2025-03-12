using AiCodeEditor.Cli.Services;
using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;

namespace AiCodeEditor.Cli.Commands
{
    [Command("search", Description = "Create a vector database from current codebase and search it")]
    public class SearchCodeCommand : ICommand
    {
        [CommandOption("query", 'q', Description = "Search query")]
        public required string Query { get; init; }

        [CommandOption("chunk-size", 's', Description = "Maximum size of each chunk in characters")]
        public int ChunkSize { get; init; } = 2048;

        [CommandOption("limit", 'l', Description = "Maximum number of results to return")]
        public int Limit { get; init; } = 5;

        [CommandOption("threshold", 't', Description = "Minimum similarity score (0-1)")]
        public float Threshold { get; init; } = 0.3f;

        [CommandOption("ollama-endpoint", Description = "Ollama API endpoint")]
        public string OllamaEndpoint { get; init; } = "http://localhost:11434";

        [CommandOption("ollama-model", Description = "Ollama model to use")]
        public string OllamaModel { get; init; } = "nomic-embed-text:latest";

        [CommandOption("qdrant-host", Description = "Qdrant host")]
        public string QdrantHost { get; init; } = "localhost";

        [CommandOption("qdrant-port", Description = "Qdrant port")]
        public int QdrantPort { get; init; } = 6334;

        public async ValueTask ExecuteAsync(IConsole console)
        {
            try
            {
                // Create services with a unique collection name
                string collectionName = Guid.NewGuid().ToString("N");
                var chunkingService = new CodebaseChunkingService(ChunkSize);
                var embeddingService = new OllamaEmbeddingService(OllamaEndpoint, OllamaModel);
                var qdrantService = new QdrantService(collectionName, QdrantHost, QdrantPort);
                var searchService = new CodeSearchService(embeddingService, qdrantService);

                // Initialize collection
                console.Output.WriteLine($"Initializing collection '{collectionName}'...");
                await qdrantService.InitializeCollectionAsync();

                // Chunk the codebase
                console.Output.WriteLine("Chunking codebase...");
                var chunks = await chunkingService.ChunkCodebaseAsync(Directory.GetCurrentDirectory());
                console.Output.WriteLine($"Found {chunks.Count} chunks");

                // Store chunks with embeddings
                console.Output.WriteLine("Generating embeddings and storing chunks...");
                for (int i = 1; i < chunks.Count + 1; i++)
                {
                    var chunk = chunks[i - 1];
                    var embedding = await embeddingService.GetEmbeddingAsync(chunk.Content);
                    var payload = new Dictionary<string, string>
                    {
                        { "file_path", chunk.FilePath },
                        { "start_line", chunk.StartLine.ToString() },
                        { "end_line", chunk.EndLine.ToString() },
                        { "language", chunk.Language },
                        { "content", chunk.Content }
                    };

                    await qdrantService.UpsertVectorAsync((ulong)i, embedding, payload);
                    
                    if (i % 10 == 0) // Progress indicator
                    {
                        console.Output.Write(".");
                    }
                }
                console.Output.WriteLine("\nStorage complete!");

                // Search
                console.Output.WriteLine($"\nSearching for: {Query}");
                var results = await searchService.SearchGroupedByFileAsync(Query, Limit, Threshold);

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