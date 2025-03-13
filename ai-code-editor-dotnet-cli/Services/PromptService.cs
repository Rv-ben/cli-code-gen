using Microsoft.SemanticKernel;
using Microsoft.SemanticKernel.Functions;

namespace AiCodeEditor.Cli.Services
{
    public class PromptService
    {
        private readonly Kernel _kernel;
        private readonly Dictionary<string, KernelFunction> _functions;

        public PromptService(Kernel kernel)
        {
            _kernel = kernel;
            _functions = LoadPromptFunctions();
        }

        private Dictionary<string, KernelFunction> LoadPromptFunctions()
        {
            var promptsPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "Prompts");
            var functions = new Dictionary<string, KernelFunction>();
            
            // Import all prompt functions from YAML files
            functions["SearchContextualized"] = _kernel.CreateFunctionFromPrompt(Path.Combine(promptsPath, "SearchContextualized", "skprompt.yaml"));
            functions["RefactoringPlan"] = _kernel.CreateFunctionFromPrompt(Path.Combine(promptsPath, "RefactoringPlan", "skprompt.yaml"));
            functions["CodeExplanation"] = _kernel.CreateFunctionFromPrompt(Path.Combine(promptsPath, "CodeExplanation", "skprompt.yaml"));
            functions["FindBug"] = _kernel.CreateFunctionFromPrompt(Path.Combine(promptsPath, "FindBug", "skprompt.yaml"));
            
            return functions;
        }

        public async Task<string> SearchContextualizedQueryAsync(string searchQuery, string codeContext)
        {
            var arguments = new KernelArguments
            {
                ["searchQuery"] = searchQuery,
                ["codeContext"] = codeContext
            };
            
            var result = await _functions["SearchContextualized"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }

        public async Task<string> GetRefactoringPlanAsync(string code, string language)
        {
            var arguments = new KernelArguments
            {
                ["code"] = code,
                ["language"] = language
            };
            
            var result = await _functions["RefactoringPlan"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }

        public async Task<string> GetCodeExplanationAsync(string code, string language)
        {   
            var arguments = new KernelArguments
            {
                ["code"] = code,
                ["language"] = language
            };
            
            var result = await _functions["CodeExplanation"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }
        
        public async Task<string> FindBugAsync(string code, string additionalContext)
        {
            var arguments = new KernelArguments
            {
                ["code"] = code,
                ["additionalContext"] = additionalContext
            };
            
            var result = await _functions["FindBug"].InvokeAsync(_kernel, arguments);
            return result.GetValue<string>();
        }
    }
} 