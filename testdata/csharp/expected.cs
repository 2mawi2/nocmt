using System;
using System.Collections.Generic;

/*
 * This is a sample C# program
 * with multi-line comments
 */
 
#pragma warning disable CS1591  // Missing XML documentation
#region Main Program

public class Program
{
    private const string VERSION = "1.0.0"; /* Version number */
    
    /// <summary>
    /// XML documentation comment
    /// </summary>
    public static void Main(string[] args)
    {
        Console.WriteLine("Hello, World!");
        
        /* Process arguments */
        #pragma warning disable CS0168  // Unused variable warning
        int unusedVar;
        #pragma warning restore CS0168
        
        foreach (var arg in args)
        {
            Console.WriteLine($"Argument: {arg}");
        }
    }
}

#endregion // End of main program
