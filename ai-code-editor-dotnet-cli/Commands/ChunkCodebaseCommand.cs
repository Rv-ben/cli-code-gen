using AiCodeEditor.Cli.Services;
using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;

namespace AiCodeEditor.Cli.Commands
{
    [Command("chunk", Description = "Chunk the codebase in the current directory")]
    public class ChunkCodebaseCommand : ICommand
    {
        [CommandOption("max-size", 's', Description = "Maximum size of each chunk in characters", IsRequired = false)]
        public int MaxChunkSize { get; init; } = 2048;

        [CommandOption("verbose", 'v', Description = "Show detailed output")]
        public bool Verbose { get; init; } = false;

        public async ValueTask ExecuteAsync(IConsole console)
        {
            var currentDirectory = Directory.GetCurrentDirectory();
            var service = new CodebaseChunkingService(MaxChunkSize);

            console.Output.WriteLine($"Chunking codebase in: {currentDirectory}");
            console.Output.WriteLine($"Max chunk size: {MaxChunkSize} characters");

            try
            {
                var chunks = await service.ChunkCodebaseAsync(currentDirectory);
                console.Output.WriteLine($"\nFound {chunks.Count} chunks across the codebase.");

                if (Verbose)
                {
                    foreach (var chunk in chunks)
                    {
                        console.Output.WriteLine("\n-------------------");
                        console.Output.WriteLine($"File: {chunk.FilePath}");
                        console.Output.WriteLine($"Lines: {chunk.StartLine}-{chunk.EndLine}");
                        console.Output.WriteLine($"Language: {chunk.Language}");
                        console.Output.WriteLine($"Size: {chunk.Content.Length} characters");
                        console.Output.WriteLine("Content preview (first 100 chars):");
                        console.Output.WriteLine(chunk.Content.Length > 100 
                            ? chunk.Content[..100] + "..."
                            : chunk.Content);
                    }
                }
            }
            catch (Exception ex)
            {
                console.Error.WriteLine($"Error processing codebase: {ex.Message}");
                throw;
            }
        }
    }
} 