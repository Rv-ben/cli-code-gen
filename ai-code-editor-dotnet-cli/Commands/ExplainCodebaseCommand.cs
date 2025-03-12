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
        [CommandParameter(0, Description = "Optional query to focus the explanation on specific parts of the codebase", IsRequired = false)]
        public string? Query { get; init; }

        private readonly PromptService _promptService;
        private readonly CodeSearchPlugin _codeSearchPlugin;

        public ExplainCodebaseCommand(PromptService promptService, CodeSearchPlugin codeSearchPlugin)
        {
            _promptService = promptService;
            _codeSearchPlugin = codeSearchPlugin;
        }

        public async ValueTask ExecuteAsync(IConsole console)
        {
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

            // Construct the prompt based on whether we have a specific query or not
            string prompt;
            if (string.IsNullOrEmpty(Query))
            {
                prompt = "Please analyze this code and provide a clear, high-level explanation of its purpose, " +
                        "main components, and how they work together. Focus on the key functionality and architecture.";
            }
            else
            {
                prompt = $"Please explain the following code, focusing on {Query}. " +
                        "Provide a clear explanation of how it works and its purpose in the codebase.";
            }

            try
            {
                var explanation = await _promptService.GetCodeExplanationAsync(prompt, relevantCode);
                await console.Output.WriteLineAsync(explanation);
            }
            catch (Exception ex)
            {
                await console.Output.WriteLineAsync($"Error getting explanation: {ex.Message}");
            }
        }
    }
} 