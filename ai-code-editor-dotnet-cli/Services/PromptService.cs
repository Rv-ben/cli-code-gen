using Microsoft.SemanticKernel;
using Microsoft.SemanticKernel.Connectors.Ollama;
using Microsoft.SemanticKernel.Functions;
using System.ComponentModel;

namespace AiCodeEditor.Cli.Services
{
    public class PromptService
    {
        private readonly Kernel _kernel;
        private readonly Dictionary<string, KernelPlugin> _functions;

        public PromptService(Kernel kernel)
        {
            _kernel = kernel;
            _functions = LoadPromptFunctions();
        }

        private Dictionary<string, KernelPlugin> LoadPromptFunctions()
        {
            var promptsPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "Prompts");
            var functions = new Dictionary<string, KernelPlugin>();

            var codeExplanationPlugin = Path.Combine(promptsPath, "CodeExplanation");
            var searchContextualizedPlugin = Path.Combine(promptsPath, "SearchContextualized");

            Console.WriteLine("Loading prompt functions");
            
            // Import all prompt functions from YAML files
            functions["CodeExplanation"] = _kernel.ImportPluginFromPromptDirectory(codeExplanationPlugin);
            functions["SearchContextualized"] = _kernel.ImportPluginFromPromptDirectory(searchContextualizedPlugin);

            Console.WriteLine("Loaded prompt functions");
            Console.WriteLine(functions["CodeExplanation"].FunctionCount);
            
            return functions;
        }

        [KernelFunction, Description("Get an explanation of the provided code")]
        public async Task<string> GetCodeExplanationAsync(
            [Description("The code to explain")] string code,
            [Description("The programming language of the code")] string language)
        {   
            var arguments = new KernelArguments
            {
                ["code"] = code,
                ["language"] = language
            };
            
            var result = await _functions["CodeExplanation"]["ExplainCode"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }

        [KernelFunction, Description("Generate a PlantUML diagram from code")]
        public async Task<string> GetPlantUMLAsync(
            [Description("The query describing what to diagram")] string query,
            [Description("The code to generate a diagram from")] string code,
            [Description("The programming language of the code")] string language)
        {
            var arguments = new KernelArguments
            {
                ["query"] = query,
                ["code"] = code,
                ["language"] = language
            };

            var result = await _functions["CodeExplanation"]["MakePlantUML"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }
        
        [KernelFunction, Description("Find potential bugs in the code")]
        public async Task<string> FindBugAsync(
            [Description("The code to analyze for bugs")] string code,
            [Description("Additional context about the code")] string additionalContext)
        {
            var arguments = new KernelArguments
            {
                ["code"] = code,
                ["additionalContext"] = additionalContext
            };
            
            var result = await _functions["FindBug"]["FindBug"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }

        [KernelFunction, Description("Generate enhanced search queries based on code context")]
        public async Task<string> GetEnhancedSearchQueryAsync(
            [Description("The original search query")] string searchQuery,
            [Description("The code context to use for enhancement")] string codeContext,
            [Description("Number of enhanced queries to generate")] int queryCount)
        {
            var arguments = new KernelArguments
            {
                ["searchQuery"] = searchQuery,
                ["codeContext"] = codeContext,
                ["queryCount"] = queryCount
            };
            
            var result = await _functions["SearchContextualized"]["EnhanceSearchQuery"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }

        [KernelFunction, Description("Generate a PlantUML diagram from code")]
        public async Task<string> GetPlantUMLV2Async(
            [Description("The query describing what to diagram")] string query,
            [Description("The code to generate a diagram from")] string code,
            [Description("The programming language of the code")] string language)
        {

            var ollamaSettings = new OllamaPromptExecutionSettings
            {
                FunctionChoiceBehavior = FunctionChoiceBehavior.Auto()
            };
            
            var arguments = new KernelArguments
            {
                ["query"] = query,
                ["code"] = code,
                ExecutionSettings = new Dictionary<string, PromptExecutionSettings>
                {
                    ["default"] = ollamaSettings
                }
            };
            
            var result = await _functions["CodeExplanation"]["MakePlantUMLV2"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }
    }
} 