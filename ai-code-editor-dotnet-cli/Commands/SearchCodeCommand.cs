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
        private readonly CodebaseIndexingService _indexingService;
        private readonly CodeSearchService _searchService;
        private readonly AppConfig _config;

        [CommandOption("query", 'q', Description = "Search query")]
        public required string Query { get; init; }

        [CommandOption("limit", 'l', Description = "Maximum number of results to return")]
        public int Limit { get; init; } = 5;

        public SearchCodeCommand(
            CodebaseIndexingService indexingService,
            CodeSearchService searchService,
            AppConfig config)
        {
            _indexingService = indexingService;
            _searchService = searchService;
            _config = config;
        }

        public async ValueTask ExecuteAsync(IConsole console)
        {
            try
            {
                // Index the codebase
                await _indexingService.IndexCodebaseAsync(
                    Directory.GetCurrentDirectory(),
                    (message, isNewLine) => 
                    {
                        if (isNewLine)
                            console.Output.WriteLineAsync(message);
                        else
                            console.Output.WriteAsync(message);
                    }
                );

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