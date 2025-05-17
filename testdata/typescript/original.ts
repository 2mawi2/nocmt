// @ts-check
// TypeScript demo file

/// <reference path="./types.d.ts" />

// Generic interface example
interface Box<T> {
	value: T; // inline trailing
}

/* Function with generics */
function identity<T>(arg: T): T {
	// inside comment
	return arg;
}

// Async / await with comments
const fetchData = async (): Promise<string> => {
	// Fetching
	const res = await fetch("https://api.example.com"); /* inline */
	return await res.text(); // trailing
};

// @ts-ignore – must stay
const bad: number = "not-a-number"; // real comment 