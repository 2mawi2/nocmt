package processor




type CSharpProcessor struct {
	*CSharpSingleProcessor
}



func NewCSharpProcessor(preserveDirectivesFlag bool) *CSharpProcessor {
	return &CSharpProcessor{
		CSharpSingleProcessor: NewCSharpSingleProcessor(preserveDirectivesFlag),
	}
}


func isCSharpDirective(line string) bool {
	return checkCSharpDirective(line)
}
