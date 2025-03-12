using CliFx;
using System.Threading.Tasks;
using Microsoft.Extensions.DependencyInjection;
using AiCodeEditor.Cli.Services;
using AiCodeEditor.Cli.Plugins;
using AiCodeEditor.Cli.Models;
using AiCodeEditor.Cli.Commands;

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
                OllamaModel = "llama2",
                QdrantHost = "localhost",
                QdrantPort = 6334,
                QdrantCollection = Guid.NewGuid().ToString("N"),
                SearchThreshold = 0.3f,
                MaxSearchResults = 2
            };
            services.AddSingleton(config);

            // Register services
            services.AddSingleton<OllamaEmbeddingService>();
            services.AddSingleton<QdrantService>();
            services.AddSingleton<SemanticKernelService>();
            services.AddSingleton<CodeSearchPlugin>();
            services.AddSingleton<PromptService>();
            services.AddSingleton<CodeSearchService>();
            
            // Register commands and their dependencies
            services.AddTransient<SearchCodeCommand>();
            services.AddTransient<CodebaseChunkingService>();

            var serviceProvider = services.BuildServiceProvider();

            return await new CliApplicationBuilder()
                .AddCommandsFromThisAssembly()
                .UseTypeActivator(type => serviceProvider.GetRequiredService(type))
                .Build()
                .RunAsync();
        }
    }
} 