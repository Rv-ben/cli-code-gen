using AiCodeEditor.Cli.Services;
using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;
using System.Threading.Tasks;

namespace AiCodeEditor.Cli.Commands
{
    [Command("generate", Description = "Generate code based on a description")]
    public class GenerateCodeCommand : ICommand
    {
        [CommandParameter(0, Description = "Description of the code to generate")]
        public required string Description { get; init; }

        [CommandOption("api-key", 'k', Description = "OpenAI API key (not needed for Ollama)", EnvironmentVariable = "OPENAI_API_KEY")]
        public string? ApiKey { get; init; }

        [CommandOption("model", 'm', Description = "Model to use (default: gpt-3.5-turbo for OpenAI, llama2 for Ollama)")]
        public string? Model { get; init; }

        [CommandOption("use-ollama", Description = "Use Ollama instead of OpenAI")]
        public bool UseOllama { get; init; } = false;

        [CommandOption("ollama-endpoint", Description = "Ollama API endpoint (default: http://localhost:11434)")]
        public string? OllamaEndpoint { get; init; }

        [CommandOption("output", 'o', Description = "Output file path (if not specified, prints to console)")]
        public string? OutputPath { get; init; }

        public async ValueTask ExecuteAsync(IConsole console)
        {
            console.Output.WriteLine("Generating code...");
            
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
            
            var code = await service.GenerateCodeAsync(Description);
            
            if (string.IsNullOrEmpty(OutputPath))
            {
                console.Output.WriteLine(code);
            }
            else
            {
                await File.WriteAllTextAsync(OutputPath, code);
                console.Output.WriteLine($"Code written to {OutputPath}");
            }
        }
    }
} 