using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;
using AiCodeEditor.Cli.Services;
using AiCodeEditor.Cli.Plugins;

namespace AiCodeEditor.Cli.Commands
{
    [Command("bug", Description = "Find potential bugs in the codebase using AI analysis")]
    public class FindBugCommand : ICommand
    {
        [CommandOption("query", 'q', Description = "Optional query to focus the bug search on specific parts of the codebase", IsRequired = false)]
        public string? Query { get; init; }

        private readonly PromptService _promptService;
        private readonly CodeSearchPlugin _codeSearchPlugin;
        private readonly CodebaseIndexingService _codebaseIndexingService;

        public FindBugCommand(
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
            try
            {
                // Index the codebase first
                await _codebaseIndexingService.IndexCodebaseAsync(
                    Directory.GetCurrentDirectory(),
                    (message, isNewLine) => 
                    {
                        console.Output.Write(message);
                        if (isNewLine) console.Output.WriteLine();
                    }
                );

                // Search for relevant code if query is provided
                string relevantCode = "";
                string searchContext = "";
                if (!string.IsNullOrEmpty(Query))
                {
                    relevantCode = await _codeSearchPlugin.SearchCode(Query);
                    if (relevantCode == "No relevant code found.")
                    {
                        await console.Output.WriteLineAsync("Could not find relevant code matching your query.");
                        return;
                    }
                    searchContext = $"The user is specifically concerned about: {Query}";
                }
                else
                {
                    // If no query provided, get some initial code to analyze
                    relevantCode = await _codeSearchPlugin.SearchCode("main entry point program.cs");
                    searchContext = "This is the main entry point of the application. Look for potential initialization or configuration issues.";
                }

                // Find potential bugs
                var bugAnalysis = await _promptService.FindBugAsync(relevantCode, searchContext);
                await console.Output.WriteLineAsync("\nBug Analysis:");
                await console.Output.WriteLineAsync("═══════════════════════════════════════");
                await console.Output.WriteLineAsync(bugAnalysis);
            }
            catch (Exception ex)
            {
                await console.Output.WriteLineAsync($"Error analyzing code for bugs: {ex.Message}");
            }
        }
    }
} 