using System.Text.Json;
using System.Text.Json.Serialization;
using AiCodeEditor.Cli.Models;

namespace AiCodeEditor.Cli.Services
{
    public class OllamaEmbeddingService
    {
        private readonly HttpClient _httpClient;
        private readonly string _host;
        private readonly string _model;

        public OllamaEmbeddingService(AppConfig config)
        {
            _httpClient = new HttpClient();
            _host = config.OllamaHost;
            _model = config.EmbeddingModel;
        }

        public record EmbeddingRequest(string Model, string Prompt);
        
        public record EmbeddingResponse
        {
            [JsonPropertyName("embedding")]
            public float[] Embedding { get; init; } = Array.Empty<float>();
        }

        public async Task<float[]> GetEmbeddingAsync(string text)
        {
            var request = new EmbeddingRequest(_model, text);
            var response = await _httpClient.PostAsync(
                $"{_host}/api/embeddings",
                new StringContent(JsonSerializer.Serialize(request))
            );

            response.EnsureSuccessStatusCode();
            var content = await response.Content.ReadAsStringAsync();
            var result = JsonSerializer.Deserialize<EmbeddingResponse>(content);

            return result?.Embedding ?? Array.Empty<float>();
        }
    }
} 