using Microsoft.SemanticKernel;
using Microsoft.SemanticKernel.Connectors.OpenAI;
using Microsoft.SemanticKernel.Connectors.Ollama;

namespace AiCodeEditor.Cli.Services
{
    public class SemanticKernelService
    {
        private readonly Kernel _kernel;

        public SemanticKernelService(string apiKey, string modelId = "llama2", bool useOllama = false, string? ollamaEndpoint = null)
        {
            var builder = Kernel.CreateBuilder();
            
            if (useOllama)
            {
                // Default Ollama endpoint if not specified
                ollamaEndpoint ??= "http://localhost:11434";

#pragma warning disable SKEXP0070 // Type is for evaluation purposes only and is subject to change or removal in future updates. Suppress this diagnostic to proceed.
                builder.AddOllamaChatCompletion(
                    modelId: modelId,
                    endpoint: new Uri(ollamaEndpoint)
                );
#pragma warning restore SKEXP0070 // Type is for evaluation purposes only and is subject to change or removal in future updates. Suppress this diagnostic to proceed.
            }
            else
            {
                builder.AddOpenAIChatCompletion(
                    modelId: modelId,
                    apiKey: apiKey
                );
            }
            
            _kernel = builder.Build();
        }

        public async Task<string> AskAsync(string prompt)
        {
            var function = _kernel.CreateFunctionFromPrompt(
                prompt,
                new PromptExecutionSettings { }
            );

            var result = await _kernel.InvokeAsync(function);
            return result.GetValue<string>() ?? "No response generated.";
        }

        public async Task<string> GenerateCodeAsync(string prompt)
        {
            var function = _kernel.CreateFunctionFromPrompt(
                $"Generate code based on this request: {prompt}\n" +
                "Provide only the code without explanations unless specifically asked for comments.",
                new PromptExecutionSettings {  }
            );

            var result = await _kernel.InvokeAsync(function);
            return result.GetValue<string>() ?? "No code generated.";
        }
    }
} 