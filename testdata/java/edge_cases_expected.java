
public class EdgeCases {
    
    private String field1;
    
    private String field2;
    
    private String field3;
    
    private String multipleComments;
    
    /* Single line block comment */ private String sameLineBlock;
    
    private String mixed; /* block */
    
    private String specialChars;
    
    private String doubleSlash;
    
    private String blockStart;
    
    private String blockEnd;
    
    private String todoComment;
    
    private String fixmeComment;
    
    private String xxxComment;
    
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
        String local = "value";
        
        if (true) {
            System.out.println("test");
        }
        
        return;
    }
    
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
    }
    
    @Override
    public String toString() {
        return "EdgeCases";
    }
} 