using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;
using AiCodeEditor.Cli.Services;
using AiCodeEditor.Cli.Plugins;
using Microsoft.Extensions.DependencyInjection;

namespace AiCodeEditor.Cli.Commands
{
    [Command("explain", Description = "Get an AI-powered explanation of the codebase or specific parts of it")]
    public class ExplainCodebaseCommand : ICommand
    {
        [CommandOption("query", 'q', Description = "Optional query to focus the explanation on specific parts of the codebase", IsRequired = true)]
        public string? Query { get; init; }

        private readonly PromptService _promptService;
        private readonly CodeSearchPlugin _codeSearchPlugin;

        private readonly CodebaseIndexingService _codebaseIndexingService;

        public ExplainCodebaseCommand(PromptService promptService, CodeSearchPlugin codeSearchPlugin, CodebaseIndexingService codebaseIndexingService)
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

            // If query is provided, search for specific code first
            string relevantCode = "";
            if (!string.IsNullOrEmpty(Query))
            {
                relevantCode = await _codeSearchPlugin.SearchCode(Query);
                if (relevantCode == "No relevant code found.")
                {
                    await console.Output.WriteLineAsync("Could not find relevant code matching your query.");
                    return;
                }
            }

            try
            {
                var explanation = await _promptService.GetCodeExplanationAsync(relevantCode, "C#");
                await console.Output.WriteLineAsync(explanation);
            }
            catch (Exception ex)
            {
                await console.Output.WriteLineAsync($"Error getting explanation: {ex.Message}");
            }
        }
    }
} 