namespace AiCodeEditor.Cli.Models
{
    public class AppConfig
    {
        // Ollama settings
        public string OllamaHost { get; set; } = "http://localhost:11434";
        public string OllamaModel { get; set; } = "qwen2.5:3b";
        public string EmbeddingModel { get; set; } = "nomic-embed-text:latest";
        
        // Qdrant settings
        public string QdrantHost { get; set; } = "localhost";
        public int QdrantPort { get; set; } = 6334;
        public string QdrantCollection { get; set; } = Guid.NewGuid().ToString("N");
        
        // Search settings
        public float SearchThreshold { get; set; } = 0.3f;
        public int MaxSearchResults { get; set; } = 2;

        // LLM settings
        public bool UseOllama { get; set; } = true;
        public string? OpenAIKey { get; set; }
        public string OpenAIModel { get; set; } = "gpt-4-turbo-preview";
    }
} 