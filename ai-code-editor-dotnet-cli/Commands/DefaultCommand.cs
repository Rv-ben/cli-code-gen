using CliFx;
using CliFx.Attributes;
using CliFx.Infrastructure;
using System.Threading.Tasks;

namespace AiCodeEditor.Cli.Commands
{
    [Command(Description = "AI Code Editor CLI - A tool for AI-assisted coding")]
    public class DefaultCommand : ICommand
    {
        public ValueTask ExecuteAsync(IConsole console)
        {
            console.Output.WriteLine("AI Code Editor CLI");
            console.Output.WriteLine("=================");
            console.Output.WriteLine();
            console.Output.WriteLine("This tool uses Microsoft's Semantic Kernel to provide AI-assisted coding capabilities.");
            console.Output.WriteLine("It supports both OpenAI and Ollama as AI providers.");
            console.Output.WriteLine();
            console.Output.WriteLine("Available commands:");
            console.Output.WriteLine("  ask       - Ask the AI assistant a question");
            console.Output.WriteLine("  generate  - Generate code based on a description");
            console.Output.WriteLine();
            console.Output.WriteLine("Examples:");
            console.Output.WriteLine("  Using OpenAI:");
            console.Output.WriteLine("    ai-code-editor ask \"What is dependency injection?\" --api-key YOUR_API_KEY");
            console.Output.WriteLine();
            console.Output.WriteLine("  Using Ollama (local):");
            console.Output.WriteLine("    ai-code-editor ask \"What is dependency injection?\" --use-ollama --model llama2");
            console.Output.WriteLine();
            console.Output.WriteLine("Run 'ai-code-editor [command] --help' for more information on a specific command.");
            
            return default;
        }
    }
} 