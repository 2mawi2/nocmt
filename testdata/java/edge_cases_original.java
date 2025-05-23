// Edge cases test file for Java comment processing

public class EdgeCases {
    
    // Comment at start of line
    private String field1; // Comment at end of line
    
    private String field2;// Comment with no space
    
    //Comment with no space at start
    private String field3;
    
    // Multiple consecutive comment lines
    // This is line 1
    // This is line 2  
    // This is line 3
    private String multipleComments;
    
    /* Single line block comment */ private String sameLineBlock;
    
    private String mixed; /* block */ // and line comment
    
    // Comment with special characters: @#$%^&*()
    private String specialChars;
    
    // Comment with // double slashes inside
    private String doubleSlash;
    
    // Comment with /* block start inside
    private String blockStart;
    
    // Comment with */ block end inside
    private String blockEnd;
    
    // TODO: This is a TODO comment that should be removed
    private String todoComment;
    
    // FIXME: This is a FIXME comment
    private String fixmeComment;
    
    // XXX: This is an XXX comment  
    private String xxxComment;
    
    // Comment with Unicode: café, naïve, résumé
    private String unicodeComment;
    
    // @formatter:off - this should be preserved
    private String formatterOff;
    
    // @formatter:on - this should be preserved  
    private String formatterOn;
    
    // CHECKSTYLE.OFF: VariableName - this should be preserved
    private String CHECKSTYLE_OFF;
    
    // CHECKSTYLE.ON: VariableName - this should be preserved
    private String CHECKSTYLE_ON;
    
    // NOSONAR - this should be preserved
    private String noSonar;
    
    // NOCHECKSTYLE - this should be preserved
    private String noCheckstyle;
    
    // NOFOLINT - this should be preserved  
    private String noFollint;
    
    public void method() {
        // Indented comment
        String local = "value"; // End of line comment
        
        if (true) { // Comment after brace
            // Deeply nested comment
            System.out.println("test"); // Another nested comment
        } // Comment after closing brace
        
        // Comment before return
        return; // Comment after return
    }
    
    // Empty line with comment
    
    // Another comment after empty line
    
    /**
     * Javadoc that should be preserved
     */
    public void javadocMethod() {
        /* Block comment
           that should be preserved */
        System.out.println("test");
    }
    
    /*
     * Multi-line block comment
     * with asterisks
     * should be preserved
     */
    public void blockCommentMethod() {
        // But this line comment should be removed
    }
    
    // @Override annotation comment - this is NOT a directive, should be removed
    @Override
    public String toString() {
        return "EdgeCases"; // Return comment
    }
} 