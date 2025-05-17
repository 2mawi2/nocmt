// Copyright 2025 Example Corp
// Licensed under MIT

using System;
using System.Collections.Generic; // Collection classes

/*
 * This is a sample C# program
 * with multi-line comments
 */
 
#pragma warning disable CS1591  // Missing XML documentation
#region Main Program

// Main class definition
public class Program
{
    // Constants
    private const string VERSION = "1.0.0"; /* Version number */
    
    /// <summary>
    /// XML documentation comment
    /// </summary>
    public static void Main(string[] args)
    {
        // Print a greeting
        Console.WriteLine("Hello, World!");
        
        /* Process arguments */
        #pragma warning disable CS0168  // Unused variable warning
        int unusedVar;
        #pragma warning restore CS0168
        
        // Loop through arguments
        foreach (var arg in args)
        {
            // Print each argument
            Console.WriteLine($"Argument: {arg}");
        }
    }
}

#endregion // End of main program

// End of file 