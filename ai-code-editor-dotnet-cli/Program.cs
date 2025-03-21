using CliFx;
using System.Threading.Tasks;
using Microsoft.Extensions.DependencyInjection;
using AiCodeEditor.Cli.Services;
using AiCodeEditor.Cli.Plugins;
using AiCodeEditor.Cli.Models;
using AiCodeEditor.Cli.Commands;
using Microsoft.SemanticKernel;
namespace AiCodeEditor.Cli
{
    public static class Program
    {
        public static async Task<int> Main()
        {
            var services = new ServiceCollection();

            // Register configuration
            var config = new AppConfig
            {
                OllamaHost = "http://localhost:11434",
                OllamaModel = "qwen2.5:14b",
                EmbeddingModel = "nomic-embed-text:latest",
                QdrantHost = "localhost",
                QdrantPort = 6334,
                QdrantCollection = Guid.NewGuid().ToString("N"),
                SearchThreshold = 0.4f,
                MaxSearchResults = 2,
                UseOllama = true,
                OpenAIKey = null,
                OpenAIModel = "gpt-4-turbo-preview"
            };
            services.AddSingleton(config);

            // Register services
            services.AddTransient<CodebaseChunkingService>();
            services.AddTransient<CodebasePathChunkingService>();
            services.AddSingleton<OllamaEmbeddingService>();
            services.AddSingleton<QdrantService>();
            services.AddSingleton<CodeSearchPlugin>();
            services.AddSingleton<PromptService>();
            services.AddSingleton<CodeSearchService>();
            services.AddSingleton<CodebaseIndexingService>();
            services.AddSingleton<IOPlugin>();
            services.AddSingleton<Kernel>(s => {
                var builder = Kernel.CreateBuilder();
                if (config.UseOllama) {
                    builder.AddOllamaChatCompletion(
                        serviceId: "ollama",
                        modelId: config.OllamaModel,
                        endpoint: new Uri(config.OllamaHost)
                    );
                }
                else {
                    builder.AddOpenAIChatCompletion(
                        modelId: config.OpenAIModel,
                        apiKey: config.OpenAIKey
                    );
                }

                builder.Services.AddSingleton<OllamaEmbeddingService>();

                builder.Plugins.AddFromType<CodeSearchPlugin>();
                builder.Plugins.AddFromType<IOPlugin>();

                services.AddTransient<CodebaseChunkingService>();
                builder.Services.AddSingleton<CodebasePathChunkingService>();
                builder.Services.AddSingleton<OllamaEmbeddingService>();
                builder.Services.AddSingleton<QdrantService>();
                builder.Services.AddSingleton<CodeSearchPlugin>();
                builder.Services.AddSingleton<PromptService>();
                builder.Services.AddSingleton<CodeSearchService>();
                builder.Services.AddSingleton<CodebaseIndexingService>();
                builder.Services.AddSingleton<IOPlugin>();
                builder.Services.AddSingleton(config);

                return builder.Build();
            });
            
            // Register commands and their dependencies
            services.AddTransient<SearchCodeCommand>();
            services.AddTransient<ExplainCodebaseCommand>();
            services.AddTransient<FindBugCommand>();
            
            services.AddTransient<MakePlantUmlCommand>();
            services.AddTransient<MakePlantUmlV2Command>();
            
            var serviceProvider = services.BuildServiceProvider();

            return await new CliApplicationBuilder()
                .AddCommandsFromThisAssembly()
                .UseTypeActivator(type => serviceProvider.GetRequiredService(type))
                .Build()
                .RunAsync();
        }
    }
} 