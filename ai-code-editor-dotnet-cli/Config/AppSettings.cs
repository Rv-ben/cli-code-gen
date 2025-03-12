namespace AiCodeEditor.Cli.Config
{
    public static class AppSettings
    {
        public static string GetApiKey()
        {
            var apiKey = Environment.GetEnvironmentVariable("OPENAI_API_KEY");
            if (string.IsNullOrEmpty(apiKey))
            {
                throw new InvalidOperationException(
                    "OpenAI API key not found. Please set the OPENAI_API_KEY environment variable or use the --api-key option."
                );
            }
            return apiKey;
        }
        
        public static string GetOllamaEndpoint()
        {
            return Environment.GetEnvironmentVariable("OLLAMA_ENDPOINT") ?? "http://localhost:11434";
        }
    }
} 