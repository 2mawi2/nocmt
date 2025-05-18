package processor




type PythonProcessor struct {
	*PythonSingleProcessor
}



func NewPythonProcessor(preserveDirectivesFlag bool) *PythonProcessor {
	return &PythonProcessor{
		PythonSingleProcessor: NewPythonSingleProcessor(preserveDirectivesFlag),
	}
}


func checkPythonDirective(line string) bool {
	return checkPythonSingleLineDirective(line)
}
