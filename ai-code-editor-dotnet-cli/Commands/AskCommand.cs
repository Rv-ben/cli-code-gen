using AiCodeEditor.Cli.Services;
using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;
using System.Threading.Tasks;

namespace AiCodeEditor.Cli.Commands
{
    [Command("ask", Description = "Ask the AI assistant a question")]
    public class AskCommand : ICommand
    {
        [CommandParameter(0, Description = "The question to ask the AI assistant")]
        public required string Question { get; init; }

        [CommandOption("api-key", 'k', Description = "OpenAI API key (not needed for Ollama)", EnvironmentVariable = "OPENAI_API_KEY")]
        public string? ApiKey { get; init; }

        [CommandOption("model", 'm', Description = "Model to use (default: gpt-3.5-turbo for OpenAI, llama2 for Ollama)")]
        public string? Model { get; init; }

        [CommandOption("use-ollama", Description = "Use Ollama instead of OpenAI")]
        public bool UseOllama { get; init; } = false;

        [CommandOption("ollama-endpoint", Description = "Ollama API endpoint (default: http://localhost:11434)")]
        public string? OllamaEndpoint { get; init; }

        public async ValueTask ExecuteAsync(IConsole console)
        {
            console.Output.WriteLine("Thinking...");
            
            string modelId = Model ?? (UseOllama ? "llama2" : "gpt-3.5-turbo");
            
            if (!UseOllama && string.IsNullOrEmpty(ApiKey))
            {
                console.Output.WriteLine("Error: API key is required when not using Ollama.");
                return;
            }
            
            var service = new SemanticKernelService(
                apiKey: ApiKey ?? string.Empty, 
                modelId: modelId,
                useOllama: UseOllama,
                ollamaEndpoint: OllamaEndpoint
            );
            
            var response = await service.AskAsync(Question);
            
            console.Output.WriteLine(response);
        }
    }
} 