namespace AiCodeEditor.Cli.Models
{
    public class AppConfig
    {
        public string OllamaHost { get; set; } = "http://localhost:11434";
        public string OllamaModel { get; set; } = "qwen2.5:3b";
        public string EmbeddingModel { get; set; } = "nomic-embed-text:latest";
        public string QdrantHost { get; set; } = "localhost";
        public int QdrantPort { get; set; } = 6334;
        public string QdrantCollection { get; set; } = Guid.NewGuid().ToString("N");
        public float SearchThreshold { get; set; } = 0.3f;
        public int MaxSearchResults { get; set; } = 2;
    }
} 