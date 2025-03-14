using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;
using AiCodeEditor.Cli.Services;
using AiCodeEditor.Cli.Plugins;
using System.Text.Json;

namespace AiCodeEditor.Cli.Commands
{
    [Command("search-context", Description = "Search code with contextual query enhancement")]
    public class SearchContextualizedCommand : ICommand
    {
        [CommandOption("query", 'q', Description = "Search query to be enhanced with context", IsRequired = true)]
        public string? Query { get; init; }

        private readonly PromptService _promptService;
        private readonly CodeSearchPlugin _codeSearchPlugin;
        private readonly CodebaseIndexingService _codebaseIndexingService;

        public SearchContextualizedCommand(
            PromptService promptService, 
            CodeSearchPlugin codeSearchPlugin, 
            CodebaseIndexingService codebaseIndexingService)
        {
            _promptService = promptService;
            _codeSearchPlugin = codeSearchPlugin;
            _codebaseIndexingService = codebaseIndexingService;
        }

        public async ValueTask ExecuteAsync(IConsole console)
        {
            await console.Output.WriteLineAsync($"Original query: {Query}");
            var currentDirectory = Directory.GetCurrentDirectory();
            await console.Output.WriteLineAsync($"Current directory: {currentDirectory}");

            await _codebaseIndexingService.IndexCodebaseAsync(
                currentDirectory,
                (message, isNewLine) => 
                {
                    console.Output.Write(message);
                    if (isNewLine) console.Output.WriteLine();
                }
            );

            // Get initial context
            string codeContext = await _codeSearchPlugin.SearchCodeFiles(Query, 1);
            if (codeContext == "No relevant files found.")
            {
                await console.Output.WriteLineAsync("Could not find initial context for your query.");
                return;
            }

            try
            {
                // Generate enhanced search query using the context
                var enhancedQuery = await _promptService.GetEnhancedSearchQueryAsync(Query, codeContext, 3);
                await console.Output.WriteLineAsync($"\nEnhanced query: {enhancedQuery}");

                // Parse the enhanced query
                var enhancedQueries = JsonSerializer.Deserialize<List<string>>(enhancedQuery);
                if (enhancedQueries == null)
                {
                    await console.Output.WriteLineAsync("Could not parse enhanced query.");
                    return;
                }

                var foundFilePaths = new List<string>();

                foreach (var query in enhancedQueries)
                {
                    await console.Output.WriteLineAsync($"\nSearching with query: {query}");
                    var filePaths = await _codeSearchPlugin.SearchFilePaths(query, 1, 0.5f, foundFilePaths);
                    foundFilePaths.AddRange(filePaths);
                    await console.Output.WriteLineAsync("\nSearch results:");
                    await console.Output.WriteLineAsync(string.Join("\n", filePaths));
                }

                // Search with enhanced query
                // var searchResults = await _codeSearchPlugin.SearchCode(enhancedQuery);
                // await console.Output.WriteLineAsync("\nSearch results:");
                // await console.Output.WriteLineAsync(searchResults);
            }
            catch (Exception ex)
            {
                await console.Output.WriteLineAsync($"Error performing contextualized search: {ex.Message}");
            }
        }
    }
}
