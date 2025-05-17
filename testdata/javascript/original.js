/* @license
 * Example JS library
 * Copyright 2025
 */

 // @flow
// @jsx React.createElement

// First comment
/* Block comment before code */
function main() {            // trailing comment on signature
	// Code line comment
	console.log("Hello");     /* inline block */ // twin trailing
	
	const url = "https://example.com/#hash"; // URL inside string
	const tmpl = `Template string with // pseudo-comment
	and /* block */ markers inside`;          // real trailing comment

	// Directive that must stay
	//# sourceMappingURL=app.js.map
}

/* A simple block comment at the end. */ 