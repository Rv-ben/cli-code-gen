using Microsoft.SemanticKernel;

namespace AiCodeEditor.Cli.Services
{
    public class PromptService
    {
        private readonly SemanticKernelService _kernel;
        private readonly Dictionary<string, PromptTemplate> _templates;

        public record PromptTemplate
        {
            public required string Name { get; init; }
            public required string Description { get; init; }
            public required string Template { get; init; }
        }

        public PromptService(string? apiKey = null, string? modelId = null, bool useOllama = true, string? ollamaEndpoint = null)
        {
            _kernel = new SemanticKernelService(
                apiKey ?? string.Empty,
                modelId ?? "llama2",
                useOllama,
                ollamaEndpoint
            );
        }

        public async Task<string> SearchContextualizedQueryAsync(string searchQuery, string codeContext)
        {
            var prompt = @$"
                Given this code context and search query, generate a more specific search query:

                Original Query: {searchQuery}

                Code Context:
                {codeContext}

                Generate a detailed search query that will find the most relevant code based on this context.
                Return only the search query without explanation.
            ";
            return await _kernel.AskAsync(prompt);
        }

        public async Task<string> GetRefactoringPlanAsync(string code, string language)
        {
            var prompt = @$"
                Create a step-by-step refactoring plan for this {language} code:

                {code}

                Plan should:
                1. Identify areas needing improvement
                2. Suggest specific refactoring techniques
                3. Maintain functionality while improving code quality
                4. List steps in order of priority
            ";
            return await _kernel.AskAsync(prompt);
        }

        public async Task<string> GetCodeExplanationAsync(string code, string language)
        {
            var prompt = @$"
                You are a senior software developer. Explain this {language} code concisely:

                {code}
            ";
            return await _kernel.AskAsync(prompt);
        }
        
    }
} 