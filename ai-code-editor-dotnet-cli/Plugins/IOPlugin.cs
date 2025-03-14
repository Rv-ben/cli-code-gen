using Microsoft.SemanticKernel;
using System.ComponentModel;

public static class IOPlugin
{

    [KernelFunction, Description("Read a file")]
    public static async Task<string> ReadFileAsync(string filePath)
    {
        var content = await File.ReadAllTextAsync(filePath);
        var absolutePath = Path.GetFullPath(filePath);
        return $"=== FILE PATH ===\n{absolutePath}\n=== CONTENT ===\n{content}";
    }

    [KernelFunction, Description("Write a file")]
    public static async Task WriteFileAsync(string filePath, string content)
    {
        await File.WriteAllTextAsync(filePath, content);
    }

    public static async Task<string> ReadFilesAsync(List<string> foundFilePaths)
    {
        var files = new System.Text.StringBuilder();
        foreach (var filePath in foundFilePaths)
        {
            if (files.Length > 0)
            {
                files.AppendLine("═══════════════════════════════════════");
                files.AppendLine();
            }
            var content = await ReadFileAsync(filePath);
            files.AppendLine(content);
        }
        return files.ToString();
    }
}