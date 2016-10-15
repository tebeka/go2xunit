# Regression Test Files

The file name should be in the following convention:
    
    <kind>-<desc>.out

Where `kind` is either `gotest` or `gocheck`.
Example: `gocheck-nofiles.out`

If the there are errors in the output (failed tests ...) add `-fail` suffix to
the file name.
Example: `gotest-fail.out`

Each of these files should have corresponding XMLs in `xml/xunit` and `xml/xunit.net/` which has the same file name with `.xml` suffix.
Example: `xml/xunit/gotest-fail.out.xml`
