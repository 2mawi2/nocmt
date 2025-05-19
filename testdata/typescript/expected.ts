/* @license
 * Example JS library
 * Copyright 2025
 */

 // @flow
// @jsx React.createElement

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
		this.logger.log(message);
		return message;
	}
}

function main(): void {
	const logger: Logger = {
		log: (msg: string) => {
			console.log("LOG:", msg); /* inline block */
		}
	};
	
	const url: string = "https://example.com/#hash";
	const tmpl: string = `Template string with // pseudo-comment
	and /* block */ markers inside`;

	const greeter: Greeter = new Greeter("Hello", logger);
	const names: string[] = ["Alice", "Bob", "Charlie"];
	names.forEach((name) => {
		greeter.greet(name);
	});

	//# sourceMappingURL=app.js.map
}

/* A simple block comment at the end. */ 