using System;
using System.Diagnostics.CodeAnalysis;
using System.IO;
using System.Net;
using System.Runtime.InteropServices;
using static System.Console;

StreamWriter? fileWriter = null;

void WriteLineToConsoleAndFile(string? value = null)
{
    WriteLine(value);
    fileWriter?.WriteLine(value);
}

// Variant of https://github.com/dotnet/core/tree/main/samples/dotnet-runtimeinfo
// Ascii text: https://ascii.co.uk/text (Univers font)

try
{
    fileWriter = new StreamWriter("output.txt");

    WriteLineToConsoleAndFile("""
         42                                                    
         42              ,d                             ,d     
         42              42                             42     
 ,adPPYb,42  ,adPPYba, MM42MMM 8b,dPPYba,   ,adPPYba, MM42MMM  
a8"    `Y42 a8"     "8a  42    42P'   `"8a a8P_____42   42     
8b       42 8b       d8  42    42       42 8PP!!!!!!!   42     
"8a,   ,d42 "8a,   ,a8"  42,   42       42 "8b,   ,aa   42,    
 `"8bbdP"Y8  `"YbbdP"'   "Y428 42       42  `"Ybbd8"'   "Y428  

""");

const double Mebi = 1024 * 1024;
const double Gibi = Mebi * 1024;
GCMemoryInfo gcInfo = GC.GetGCMemoryInfo();
long totalMemoryBytes = gcInfo.TotalAvailableMemoryBytes;

// OS and .NET information
WriteLineToConsoleAndFile($"{nameof(RuntimeInformation.OSArchitecture)}: {RuntimeInformation.OSArchitecture}");
WriteLineToConsoleAndFile($"{nameof(RuntimeInformation.OSDescription)}: {RuntimeInformation.OSDescription}");
WriteLineToConsoleAndFile($"{nameof(RuntimeInformation.FrameworkDescription)}: {RuntimeInformation.FrameworkDescription}");
WriteLineToConsoleAndFile();

// Environment information
WriteLineToConsoleAndFile($"{nameof(Environment.UserName)}: {Environment.UserName}");
WriteLineToConsoleAndFile($"HostName : {Dns.GetHostName()}");
WriteLineToConsoleAndFile();

// Hardware information
WriteLineToConsoleAndFile($"{nameof(Environment.ProcessorCount)}: {Environment.ProcessorCount}");
WriteLineToConsoleAndFile($"{nameof(GCMemoryInfo.TotalAvailableMemoryBytes)}: {totalMemoryBytes} ({GetInBestUnit(totalMemoryBytes)})");

string[] memoryLimitPaths = new string[] 
{
    "/sys/fs/cgroup/memory.max",
    "/sys/fs/cgroup/memory.high",
    "/sys/fs/cgroup/memory.low",
    "/sys/fs/cgroup/memory/memory.limit_in_bytes",
};

string[] currentMemoryPaths = new string[] 
{
    "/sys/fs/cgroup/memory.current",
    "/sys/fs/cgroup/memory/memory.usage_in_bytes",
};

// cgroup information
if (OperatingSystem.IsLinux() &&
    GetBestValue(memoryLimitPaths, out long memoryLimit, out string? bestMemoryLimitPath) &&
    memoryLimit > 0)
{
    // get memory cgroup information
    GetBestValue(currentMemoryPaths, out long currentMemory, out string? memoryPath);

    WriteLineToConsoleAndFile($"cgroup memory constraint: {bestMemoryLimitPath}");
    WriteLineToConsoleAndFile($"cgroup memory limit: {memoryLimit} ({GetInBestUnit(memoryLimit)})");
    WriteLineToConsoleAndFile($"cgroup memory usage: {currentMemory} ({GetInBestUnit(currentMemory)})");
    WriteLineToConsoleAndFile($"GC Hard limit %: {(double)totalMemoryBytes/memoryLimit * 100:N0}");
}

string GetInBestUnit(long size)
{
    if (size < Mebi)
    {
        return $"{size} bytes";
    }
    else if (size < Gibi)
    {
        double mebibytes = size / Mebi;
        return $"{mebibytes:F} MiB";
    }
    else
    {
        double gibibytes = size / Gibi;
        return $"{gibibytes:F} GiB";
    }
}

bool GetBestValue(string[] paths, out long limit, [NotNullWhen(true)] out string? bestPath)
{
    foreach (string path in paths)
    {
        if (Path.Exists(path) &&
            long.TryParse(File.ReadAllText(path), out limit))
        {
            bestPath = path;
            return true;
        }
    }

    bestPath = null;
    limit = 0;
    return false;
}
}
finally
{
    fileWriter?.Close();
}

