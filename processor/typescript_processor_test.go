package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeScriptStripComments(t *testing.T) {
	t.Run("FileBased", func(t *testing.T) {
		processor := NewTypeScriptProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../testdata/typescript/original.ts", "../testdata/typescript/expected.ts")
	})

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "strip comments in TypeScript syntax",
			input: `// TypeScript interface
interface User {
	name: string; // user name
	age: number; /* user age */
}

// TypeScript class with type annotations
class UserService {
	private users: User[]; // array of users
	
	constructor() {
		this.users = []; // initialize empty array
	}
	
	addUser(user: User): void { // add user method
		this.users.push(user);
	}
}`,
			expected: `interface User {
	name: string;
	age: number;
}

class UserService {
	private users: User[];
	
	constructor() {
		this.users = [];
	}
	
	addUser(user: User): void {
		this.users.push(user);
	}
}`,
		},
		{
			name: "strip comments in generic types",
			input: `function identity<T>(arg: T): T {
	// This is a generic function
	return arg;
}

// Generic interface
interface GenericInterface<T> {
	value: T; // Generic value
	method<U>(arg: U /* arg comment */): U; // Generic method
}

// Generic with constraints
interface WithLength {
	length: number;
}

function loggingIdentity<T extends WithLength>(arg: T): T {
	/* Log the length */
	console.log(arg.length);
	return arg;
}`,
			expected: `function identity<T>(arg: T): T {
	return arg;
}

interface GenericInterface<T> {
	value: T;
	method<U>(arg: U): U;
}

interface WithLength {
	length: number;
}

function loggingIdentity<T extends WithLength>(arg: T): T {
	console.log(arg.length);
	return arg;
}`,
		},
		{
			name: "strip comments in async/await syntax",
			input: `// Async function example
async function fetchData(): Promise<string> {
	// Fetching data
	const response = await fetch('https://api.example.com'); // HTTP request
	/* Processing response */
	const data = await response.text();
	return data; // Return the data
}

// Async arrow function
const fetchUser = async (id: number): Promise<User> => {
	// Fetching user
	return await fetchData(); // Placeholder
}`,
			expected: `async function fetchData(): Promise<string> {
	const response = await fetch('https://api.example.com');
	const data = await response.text();
	return data;
}

const fetchUser = async (id: number): Promise<User> => {
	return await fetchData();
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewTypeScriptProcessor(false)
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			expectedNormalized := normalizeText(tt.expected)
			assert.Equal(t, expectedNormalized, result)
		})
	}
}

func TestTypeScriptProcessorGetLanguageName(t *testing.T) {
	processor := NewTypeScriptProcessor(false)
	assert.Equal(t, "typescript", processor.GetLanguageName())
}

func TestTypeScriptProcessorPreserveDirectives(t *testing.T) {
	processorWithDirectives := NewTypeScriptProcessor(true)
	processorWithoutDirectives := NewTypeScriptProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}
