using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;
using AiCodeEditor.Cli.Services;
using AiCodeEditor.Cli.Plugins;
using Microsoft.Extensions.DependencyInjection;
using System.Text.Json;

namespace AiCodeEditor.Cli.Commands
{
    [Command("plantuml", Description = "Generate PlantUML diagram from code")]
    public class MakePlantUmlCommand : ICommand
    {
        [CommandOption("query", 'q', Description = "Query to find code to generate diagram for", IsRequired = true)]
        public string? Query { get; init; }

        private readonly PromptService _promptService;
        private readonly CodeSearchPlugin _codeSearchPlugin;
        private readonly CodebaseIndexingService _codebaseIndexingService;

        public MakePlantUmlCommand(PromptService promptService, CodeSearchPlugin codeSearchPlugin, CodebaseIndexingService codebaseIndexingService)
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
            var foundFilePaths = await _codeSearchPlugin.SearchFilePaths(Query, 1);
            if (foundFilePaths.Count == 0)
            {
                await console.Output.WriteLineAsync("Could not find initial context for your query.");
                return;
            }

            await console.Output.WriteLineAsync($"Found files: {string.Join(", ", foundFilePaths)}");

            string builder = "";
            foreach (var filePath in foundFilePaths)
            {
                builder += await IOPlugin.ReadFileAsync(filePath);
            }

            try
            {
                // Generate enhanced search query using the context
                var enhancedQuery = await _promptService.GetEnhancedSearchQueryAsync(Query, builder, 2);
                await console.Output.WriteLineAsync($"\nEnhanced query: {enhancedQuery}");

                // Parse the enhanced query
                var enhancedQueries = JsonSerializer.Deserialize<List<string>>(enhancedQuery);
                if (enhancedQueries == null)
                {
                    await console.Output.WriteLineAsync("Could not parse enhanced query.");
                    return;
                }

                foreach (var query in enhancedQueries)
                {
                    await console.Output.WriteLineAsync($"\nSearching with query: {query}");
                    var filePaths = await _codeSearchPlugin.SearchFilePathsUsingCodeContext(query, 3, 0.5f, foundFilePaths);
                    var topResult = filePaths.FirstOrDefault();
                    if (topResult != null)
                    {
                        foundFilePaths.Add(topResult);
                    }
                    await console.Output.WriteLineAsync("\nSearch results:");
                    await console.Output.WriteLineAsync(string.Join("\n", filePaths));
                }
            }
            catch (Exception ex)
            {
                await console.Output.WriteLineAsync($"Error performing contextualized search: {ex.Message}");
            }
            
            // Get all the files in the foundFilePaths
            var allFiles = await IOPlugin.ReadFilesAsync(foundFilePaths);

            // generate a plant uml diagram of the codebase
            var plantUml = await _promptService.GetPlantUMLAsync(Query, allFiles, "C#");
            await console.Output.WriteLineAsync("\nPlantUML diagram:");
            await console.Output.WriteLineAsync(plantUml);
        }
    }
}
