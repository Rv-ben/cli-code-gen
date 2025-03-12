using Qdrant.Client;
using Qdrant.Client.Grpc;

namespace AiCodeEditor.Cli.Services
{
    public class QdrantService
    {
        private readonly QdrantClient _client;
        private readonly string _collectionName;
        private const int VectorSize = 768; // Default size for Ollama embeddings

        public QdrantService(string collectionName = "code_chunks", string host = "localhost", int port = 6334)
        {
            _client = new QdrantClient(host, port);
            _collectionName = collectionName;
        }

        public async Task InitializeCollectionAsync()
        {
            try
            {
                try
                {
                    var info = await _client.GetCollectionInfoAsync(_collectionName);
                }
                catch (Grpc.Core.RpcException)
                {
                    Console.WriteLine("Collection doesn't exist, creating it...");
                    // Collection doesn't exist, create it
                    await _client.CreateCollectionAsync(
                        _collectionName,
                        new VectorParams { Size = VectorSize, Distance = Distance.Cosine }
                    );
                }
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to initialize Qdrant collection: {ex.Message}", ex);
            }
        }

        public async Task UpsertVectorAsync(
            ulong id,
            float[] vector,
            Dictionary<string, string> payload)
        {
            try
            {
                var points = new[] {
                    new PointStruct
                    {
                        Id = id,
                        Vectors = vector,
                        Payload = {
                            { "file_path", payload["file_path"] },
                            { "start_line", payload["start_line"] },
                            { "end_line", payload["end_line"] },
                            { "language", payload["language"] },
                            { "content", payload["content"] }
                        },
                        
                    }
                };

                await _client.UpsertAsync(_collectionName, points);
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to upsert vector: {ex.Message}", ex);
            }
        }

        public async Task<List<ScoredPoint>> SearchAsync(
            float[] queryVector,
            int limit = 5,
            float scoreThreshold = 0.7f)
        {
            try
            {
                var searchResult = await _client.SearchAsync(
                    _collectionName,
                    queryVector,
                    limit: (uint)limit,
                    scoreThreshold: scoreThreshold
                );

                return searchResult.ToList();
            }
            catch (Exception ex)
            {
                throw new Exception($"Failed to search vectors: {ex.Message}", ex);
            }
        }
    }
} 