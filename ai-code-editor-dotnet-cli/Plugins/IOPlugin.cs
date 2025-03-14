using Microsoft.SemanticKernel;
using System.ComponentModel;

public class IOPlugin
{

    [KernelFunction, Description("Read a file")]
    public async Task<string> ReadFileAsync(
        [Description("The path of the file to read")]
        string filePath)
    {
        Console.WriteLine($"Reading file: {filePath}");
        var content = await File.ReadAllTextAsync(filePath);
        var absolutePath = Path.GetFullPath(filePath);
        return $"=== FILE PATH ===\n{absolutePath}\n=== CONTENT ===\n{content}";
    }

    [KernelFunction, Description("Write a file")]
    public async Task WriteFileAsync(
        [Description("The path of the file to write")]
        string filePath,
        [Description("The content to write to the file")]
        string content)
    {
        await File.WriteAllTextAsync(filePath, content);
    }

    [KernelFunction, Description("Read multiple files")]
    public async Task<string> ReadFilesAsync(
        [Description("The list of file paths to read")]
        List<string> foundFilePaths)
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