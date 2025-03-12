using Microsoft.SemanticKernel;
using Microsoft.SemanticKernel.Connectors.OpenAI;
using AiCodeEditor.Cli.Models;
using Microsoft.SemanticKernel.Connectors.Ollama;

namespace AiCodeEditor.Cli.Services
{
    public class SemanticKernelService
    {
        private readonly Kernel _kernel;

        private readonly OllamaPromptExecutionSettings _ollamaSettings;

        public SemanticKernelService(AppConfig config)
        {
            var builder = Kernel.CreateBuilder();

            if (config.UseOllama)
            {
                builder.AddOllamaChatCompletion(
                    serviceId: "ollama",
                    modelId: config.OllamaModel, 
                    endpoint: new Uri(config.OllamaHost)
                );

                _ollamaSettings = new OllamaPromptExecutionSettings {
                    NumPredict = 100
                };
            }
            else if (!string.IsNullOrEmpty(config.OpenAIKey))
            {
                builder.AddOpenAIChatCompletion(
                    modelId: config.OpenAIModel,
                    apiKey: config.OpenAIKey
                );
            }
            else
            {
                throw new ArgumentException("Either UseOllama must be true or OpenAIKey must be provided");
            }

            _kernel = builder.Build();
        }

        public async Task<string> AskAsync(string prompt)
        {
            try
            {
                Console.WriteLine(_ollamaSettings.ToString());
                var promptFunction = _kernel.CreateFunctionFromPrompt(prompt, _ollamaSettings);

                var args = new KernelArguments
                {
                    { "ollama", _ollamaSettings }
                };

                var result = await _kernel.InvokeAsync(promptFunction, args);
                return result.GetValue<string>() ?? string.Empty;
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to get response from LLM: {ex.Message}", ex);
            }
        }
    }
} 