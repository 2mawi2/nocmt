/* @license
 * Example JS library
 * Copyright 2025
 */

 // @flow
// @jsx React.createElement

// First comment
/* Block comment before code */

interface Logger {
	log(message: string): void;
}

class Greeter {
	private greeting: string;
	private logger: Logger;

	constructor(greeting: string, logger: Logger) {
		this.greeting = greeting;
		this.logger = logger;
	}

	greet(name: string): string {
		const message = `${this.greeting}, ${name}!`;
		this.logger.log(message); // trailing comment on method
		return message;
	}
}

function main(): void {
	// Code line comment
	const logger: Logger = {
		log: (msg: string) => {
			console.log("LOG:", msg); /* inline block */ // twin trailing
		}
	};
	
	const url: string = "https://example.com/#hash"; // URL inside string
	const tmpl: string = `Template string with // pseudo-comment
	and /* block */ markers inside`;          // real trailing comment

	const greeter: Greeter = new Greeter("Hello", logger);
	const names: string[] = ["Alice", "Bob", "Charlie"];
	names.forEach((name) => {
		greeter.greet(name);
	});

	// Directive that must stay
	//# sourceMappingURL=app.js.map
}

/* A simple block comment at the end. */ 