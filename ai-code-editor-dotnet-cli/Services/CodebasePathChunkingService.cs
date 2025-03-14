using System.Text;
using Qdrant.Client.Grpc;

namespace AiCodeEditor.Cli.Services
{
    public class CodebasePathChunkingService
    {
        private readonly int _maxChunkSize;
        private readonly string[] _allowedExtensions = new[] { 
            ".cs", ".js", ".ts", ".py", ".go", ".java", ".cpp", ".hpp", ".h", 
            ".c", ".jsx", ".tsx", ".php", ".rb", ".rs", ".swift"
        };

        private readonly string[] _ignoredPaths = new[] {
            "node_modules",
            "obj",
            ".git",
            "dist",
            "build",
            "target",
            "bin/Debug",
            "bin/Release",
            ".vs",
            ".vscode",
            "packages",
            "TestResults",
            "coverage",
            "bin",
            "obj",
            "Prompts"
        };

        public enum ChunkType
        {
            FullPath,           // Complete path
            DirectoryOnly,      // Just the directory structure
            FileNameOnly,       // Just the filename
            DirectorySegment    // Individual directory segment
        }

        public CodebasePathChunkingService(int maxChunkSize = 100)
        {
            _maxChunkSize = maxChunkSize;
        }

        public record PathChunk(
            string FilePath,
            string AbsolutePath,
            string Language,
            string DirectoryStructure,
            string Segment,
            ChunkType Type,
            int SegmentDepth
        );

        public static Google.Protobuf.Collections.MapField<string, Value> ToFieldMap(PathChunk chunk)
        {
            return new Google.Protobuf.Collections.MapField<string, Value>
            {
                { "file_path", chunk.FilePath },
                { "absolute_path", chunk.AbsolutePath },
                { "language", chunk.Language },
                { "directory_structure", chunk.DirectoryStructure },
                { "segment", chunk.Segment },
                { "chunk_type", chunk.Type.ToString() },
                { "segment_depth", chunk.SegmentDepth.ToString() }
            };
        }

        public List<PathChunk> ChunkCodebasePaths(string rootPath)
        {
            var chunks = new List<PathChunk>();
            var files = Directory.GetFiles(rootPath, "*.*", SearchOption.AllDirectories)
                .Where(file => 
                    _allowedExtensions.Contains(Path.GetExtension(file).ToLower()) &&
                    File.Exists(file) && 
                    !ShouldSkipFile(file));

            foreach (var file in files)
            {
                var relativePath = Path.GetRelativePath(rootPath, file);
                var absolutePath = Path.GetFullPath(file);
                var language = GetLanguageFromExtension(Path.GetExtension(file));
                var directoryStructure = Path.GetDirectoryName(relativePath) ?? string.Empty;
                var fileName = Path.GetFileName(file);

                // Add filename-only chunk
                chunks.Add(new PathChunk(
                    FilePath: relativePath,
                    AbsolutePath: absolutePath,
                    Language: language,
                    DirectoryStructure: directoryStructure,
                    Segment: fileName,
                    Type: ChunkType.FileNameOnly,
                    SegmentDepth: -1
                ));
            }

            return chunks;
        }

        private bool ShouldSkipFile(string filePath)
        {
            return _ignoredPaths.Any(pattern => 
                filePath.Contains(pattern, StringComparison.OrdinalIgnoreCase));
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
