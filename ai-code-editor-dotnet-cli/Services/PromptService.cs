using Microsoft.SemanticKernel;
using Microsoft.SemanticKernel.Functions;

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

        public async Task<string> GetCodeExplanationAsync(string code, string language)
        {   
            var arguments = new KernelArguments
            {
                ["code"] = code,
                ["language"] = language
            };
            
            var result = await _functions["CodeExplanation"]["ExplainCode"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }

        public async Task<string> GetPlantUMLAsync(string code, string language)
        {
            var arguments = new KernelArguments
            {
                ["code"] = code,
                ["language"] = language
            };

            var result = await _functions["CodeExplanation"]["MakePlantUML"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }
        
        public async Task<string> FindBugAsync(string code, string additionalContext)
        {
            var arguments = new KernelArguments
            {
                ["code"] = code,
                ["additionalContext"] = additionalContext
            };
            
            var result = await _functions["FindBug"]["FindBug"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }

        public async Task<string> GetEnhancedSearchQueryAsync(string searchQuery, string codeContext, int queryCount)
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
    }
} 