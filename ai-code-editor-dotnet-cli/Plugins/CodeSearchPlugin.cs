using AiCodeEditor.Cli.Services;
using Microsoft.SemanticKernel;
using System.ComponentModel;
using AiCodeEditor.Cli.Models;

namespace AiCodeEditor.Cli.Plugins
{
    public class CodeSearchPlugin
    {
        private readonly CodeSearchService _searchService;

        public CodeSearchPlugin(
            OllamaEmbeddingService embeddingService,
            QdrantService qdrantService,
            AppConfig config)
        {
            _searchService = new CodeSearchService(embeddingService, qdrantService);
        }

        [KernelFunction, Description("Search for files in the codebase")]
        public async Task<string> SearchCodeFiles(
            [Description("The search query to find relevant code")] string query, 
            [Description("The maximum number of results to return")] int maxResults = 3,
            [Description("The minimum score threshold for results")] float threshold = 0.5f)
        {
            Console.WriteLine($"\nSearching for: {query}");
            var results = await _searchService.SearchAsync(query, maxResults, threshold);
            
            Console.WriteLine($"Found {results.Count} results");
            foreach (var result in results)
            {
                Console.WriteLine($"Match: {result.FilePath} (Score: {result.Score:F2})");
            }
            
            if (!results.Any())
            {
                return "No relevant files found.";
            }

            var response = new System.Text.StringBuilder();
            foreach (var result in results)
            {
                response.AppendLine("═══════════════════════════════════════");
                response.AppendLine($"File: {result.FilePath}");
                response.AppendLine($"Match found at lines {result.StartLine}-{result.EndLine} (Score: {result.Score:F2})");
                response.AppendLine("═══════════════════════════════════════");
                response.AppendLine($"```{result.Language}");
                
                // Read and return the entire file
                var fileContent = await File.ReadAllTextAsync(result.FilePath);
                response.AppendLine(fileContent);
                response.AppendLine("```");
                response.AppendLine();
            }

            return response.ToString().TrimEnd();
        }

        public async Task<List<string>> SearchFilePaths(string query, int maxResults = 3, float threshold = 0.5f, List<string>? excludedFilePaths = null)
        {
            return await _searchService.SearchFilePathsAsync(query, maxResults, threshold, excludedFilePaths);
        }
    }
} 