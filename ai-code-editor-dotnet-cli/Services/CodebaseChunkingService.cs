using System.Text;
using Qdrant.Client.Grpc;

namespace AiCodeEditor.Cli.Services
{
    public class CodebaseChunkingService
    {
        private readonly int _maxChunkSize;
        private readonly string[] _allowedExtensions = new[] { 
            ".cs", ".js", ".ts", ".py", ".go", ".java", ".cpp", ".hpp", ".h", 
            ".c", ".jsx", ".tsx", ".php", ".rb", ".rs", ".swift" 
        };

        public CodebaseChunkingService(int maxChunkSize = 2048)
        {
            _maxChunkSize = maxChunkSize;
        }

        public record CodeChunk(
            string Content,
            string FilePath,
            int StartLine,
            int EndLine,
            string Language
        );

        public static Google.Protobuf.Collections.MapField<string, Value> ToFieldMap(CodeChunk chunk)
        {
            return new Google.Protobuf.Collections.MapField<string, Value>
            {
                { "file_path", chunk.FilePath },
                { "start_line", chunk.StartLine.ToString() },
                { "end_line", chunk.EndLine.ToString() },
                { "language", chunk.Language },
                { "content", chunk.Content }
            };
        }

        public async Task<List<CodeChunk>> ChunkCodebaseAsync(string rootPath)
        {
            var chunks = new List<CodeChunk>();
            var files = Directory.GetFiles(rootPath, "*.*", SearchOption.AllDirectories)
                .Where(file => _allowedExtensions.Contains(Path.GetExtension(file).ToLower()));

            foreach (var file in files)
            {
                if (ShouldSkipFile(file))
                    continue;

                var fileContent = await File.ReadAllTextAsync(file);
                var language = GetLanguageFromExtension(Path.GetExtension(file));
                chunks.AddRange(ChunkFile(fileContent, file, language));
            }

            return chunks;
        }

        private bool ShouldSkipFile(string filePath)
        {
            var skipPatterns = new[]
            {
                "node_modules",
                "bin",
                "obj",
                ".git",
                "dist",
                "build",
                "target"
            };

            return skipPatterns.Any(pattern => 
                filePath.Contains(pattern, StringComparison.OrdinalIgnoreCase));
        }

        private List<CodeChunk> ChunkFile(string content, string filePath, string language)
        {
            var chunks = new List<CodeChunk>();
            var lines = content.Split('\n');
            var currentChunk = new StringBuilder();
            var startLine = 0;
            var currentLine = 0;

            foreach (var line in lines)
            {
                if (currentChunk.Length + line.Length > _maxChunkSize && currentChunk.Length > 0)
                {
                    // Create chunk and reset
                    chunks.Add(new CodeChunk(
                        currentChunk.ToString().Trim(),
                        filePath,
                        startLine,
                        currentLine - 1,
                        language
                    ));

                    currentChunk.Clear();
                    startLine = currentLine;
                }

                currentChunk.AppendLine(line);
                currentLine++;
            }

            // Add remaining content
            if (currentChunk.Length > 0)
            {
                chunks.Add(new CodeChunk(
                    currentChunk.ToString().Trim(),
                    filePath,
                    startLine,
                    currentLine - 1,
                    language
                ));
            }

            return chunks;
        }

        private string GetLanguageFromExtension(string extension)
        {
            return extension.ToLower() switch
            {
                ".cs" => "csharp",
                ".js" => "javascript",
                ".ts" => "typescript",
                ".py" => "python",
                ".go" => "go",
                ".java" => "java",
                ".cpp" => "cpp",
                ".hpp" => "cpp",
                ".h" => "c",
                ".c" => "c",
                ".jsx" => "javascript",
                ".tsx" => "typescript",
                ".php" => "php",
                ".rb" => "ruby",
                ".rs" => "rust",
                ".swift" => "swift",
                _ => "text"
            };
        }
    }
} 