// @ts-check

/// <reference path="./types.d.ts" />

interface Box<T> {
	value: T;
}

function identity<T>(arg: T): T {
	return arg;
}

const fetchData = async (): Promise<string> => {
	const res = await fetch("https://api.example.com");
	return await res.text();
};

// @ts-ignore – must stay
const bad: number = "not-a-number";
