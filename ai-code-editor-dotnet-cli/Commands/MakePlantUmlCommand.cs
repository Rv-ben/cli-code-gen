using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;
using AiCodeEditor.Cli.Services;
using AiCodeEditor.Cli.Plugins;
using Microsoft.Extensions.DependencyInjection;

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
            await console.Output.WriteLineAsync($"Query: {Query}");
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

            string relevantCode = await _codeSearchPlugin.SearchCodeFiles(Query);
            if (relevantCode == "No relevant files found.")
            {
                await console.Output.WriteLineAsync("Could not find relevant code matching your query.");
                return;
            }

            try
            {
                var plantUml = await _promptService.GetPlantUMLAsync(relevantCode, "C#");
                await console.Output.WriteLineAsync(plantUml);
            }
            catch (Exception ex)
            {
                await console.Output.WriteLineAsync($"Error generating PlantUML diagram: {ex.Message}");
            }
        }
    }
}
