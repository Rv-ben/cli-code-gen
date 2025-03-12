using Microsoft.SemanticKernel;

namespace AiCodeEditor.Cli.Services
{
    public class PromptService
    {
        private readonly SemanticKernelService _kernelService;

        public record PromptTemplate
        {
            public required string Name { get; init; }
            public required string Description { get; init; }
            public required string Template { get; init; }
        }

        public PromptService(SemanticKernelService kernelService)
        {
            _kernelService = kernelService;
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
            return await _kernelService.AskAsync(prompt);
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
            return await _kernelService.AskAsync(prompt);
        }

        public async Task<string> GetCodeExplanationAsync(string code, string language)
        {
            var prompt = @$"
                You are a senior software developer.
                Keep the answer under 100 words.
                Explain this {language} code concisely:

                {code}
            ";
            return await _kernelService.AskAsync(prompt);
        }
        
        public async Task<string> FindBugAsync(string code, string additionalContext)
        {
            var prompt = @$"
                You are a senior software developer. 
                You are given a code and some additional context.
                Find the bug in the code. Keep the answer under 100 words.

                Code:
                {code}

                Additional Context:
                {additionalContext}
            ";

            return await _kernelService.AskAsync(prompt);
        }
    }
} 