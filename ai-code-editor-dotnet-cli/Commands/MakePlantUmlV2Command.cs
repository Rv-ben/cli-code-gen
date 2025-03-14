using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;
using AiCodeEditor.Cli.Services;
using AiCodeEditor.Cli.Plugins;

namespace AiCodeEditor.Cli.Commands
{
    [Command("plantuml2", Description = "Generate PlantUML diagram from code using AI-driven function calling")]
    public class MakePlantUmlV2Command : ICommand
    {
        [CommandOption("query", 'q', Description = "Query to find code to generate diagram for", IsRequired = true)]
        public string? Query { get; init; }

        private readonly CodebaseIndexingService _codebaseIndexingService;

        private readonly PromptService _promptService;

        public MakePlantUmlV2Command(PromptService promptService, CodebaseIndexingService codebaseIndexingService)
        {
            _promptService = promptService;
            _codebaseIndexingService = codebaseIndexingService;
        }

        public async ValueTask ExecuteAsync(IConsole console)
        {
            if (string.IsNullOrEmpty(Query))
            {
                await console.Output.WriteLineAsync("Query is required.");
                return;
            }

            await console.Output.WriteLineAsync($"Original query: {Query}");
            var currentDirectory = Directory.GetCurrentDirectory();
            await console.Output.WriteLineAsync($"Current directory: {currentDirectory}");

            // Index the codebase
            await console.Output.WriteLineAsync("\nIndexing codebase...");
            await _codebaseIndexingService.IndexCodebaseAsync(
                currentDirectory,
                (message, isNewLine) => 
                {
                    console.Output.Write(message);
                    if (isNewLine) console.Output.WriteLine();
                }
            );

            try
            {
                await console.Output.WriteLineAsync("\nLetting AI analyze the codebase and generate diagram...");
                // Let the AI handle the entire process of finding and processing code
                var plantUml = await _promptService.GetPlantUMLV2Async(Query, "", "");
                await console.Output.WriteLineAsync(plantUml);
            }
            catch (Exception ex)
            {
                await console.Output.WriteLineAsync($"Error generating PlantUML diagram: {ex.Message}");
                if (ex.InnerException != null)
                {
                    await console.Output.WriteLineAsync($"Inner error: {ex.InnerException.Message}");
                }
            }
        }
    }
} 