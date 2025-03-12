using AiCodeEditor.Cli.Services;
using Microsoft.SemanticKernel;
using System.ComponentModel;
using AiCodeEditor.Cli.Models;

namespace AiCodeEditor.Cli.Plugins
{
    public class CodeSearchPlugin
    {
        private readonly CodeSearchService _searchService;
        private readonly int _maxResults;
        private readonly float _threshold;

        public CodeSearchPlugin(
            OllamaEmbeddingService embeddingService,
            QdrantService qdrantService,
            AppConfig config)
        {
            _searchService = new CodeSearchService(embeddingService, qdrantService);
            _maxResults = config.MaxSearchResults;
            _threshold = config.SearchThreshold;
        }

        [KernelFunction, Description("Search code in the codebase")]
        public async Task<string> SearchCode(
            [Description("The search query to find relevant code")] string query)
        {
            var results = await _searchService.SearchAsync(query, _maxResults, _threshold);
            
            if (!results.Any())
            {
                return "No relevant code found.";
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
    }
} 