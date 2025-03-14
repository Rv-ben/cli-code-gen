using Microsoft.SemanticKernel;
using System.ComponentModel;

public static class IOPlugin
{

    [KernelFunction, Description("Read a file")]
    public static async Task<string> ReadFileAsync(string filePath)
    {
        return await File.ReadAllTextAsync(filePath);
    }

    [KernelFunction, Description("Write a file")]
    public static async Task WriteFileAsync(string filePath, string content)
    {
        await File.WriteAllTextAsync(filePath, content);
    }
}