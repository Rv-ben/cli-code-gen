using System.Text.Json;
using System.Text.Json.Serialization;

namespace AiCodeEditor.Cli.Services
{
    public class OllamaEmbeddingService
    {
        private readonly HttpClient _httpClient;
        private readonly string _baseUrl;
        private readonly string _modelName;

        public OllamaEmbeddingService(string baseUrl = "http://localhost:11434", string modelName = "llama2")
        {
            _httpClient = new HttpClient();
            _baseUrl = baseUrl.TrimEnd('/');
            _modelName = modelName;
        }

        public record EmbeddingRequest(string Model, string Prompt);
        
        public record EmbeddingResponse
        {
            [JsonPropertyName("embedding")]
            public float[] Embedding { get; init; } = Array.Empty<float>();
        }

        public async Task<float[]> GetEmbeddingAsync(string text)
        {
            var request = new EmbeddingRequest(_modelName, text);
            var response = await _httpClient.PostAsync(
                $"{_baseUrl}/api/embeddings",
                new StringContent(JsonSerializer.Serialize(request))
            );

            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            var result = JsonSerializer.Deserialize<EmbeddingResponse>(content);

            return result?.Embedding ?? Array.Empty<float>();
        }
    }
} 