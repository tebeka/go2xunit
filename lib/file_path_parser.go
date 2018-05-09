package lib

import (
	"go/build"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func fileFileForTest(suite *Suite, test *Test, pkg *build.Package, strip string) string {
	var (
		testFileContents = make(map[string]string)
		body             []byte
		testFilePath     string
		contents         string
		err              error
		ok               bool
	)
	if strip == "" {
		strip = pkg.SrcRoot + "/"
	}
	for _, testFile := range pkg.TestGoFiles {
		testFilePath = filepath.Join(pkg.Dir, testFile)
		if contents, ok = testFileContents[testFilePath]; !ok {
			if body, err = ioutil.ReadFile(testFilePath); err != nil {
				continue
			}
			contents = string(body)
			testFileContents[testFilePath] = contents
		}
		if strings.Contains(contents, test.Name) {
			return strings.Replace(testFilePath, strip, "", 1)
		}
	}

	return suite.Name
}

func tmplFuncFilePathForTest(suite *Suite, test *Test, strip string) string {
	var (
		pkg *build.Package
		err error
	)
	if pkg, err = build.Import(suite.Name, build.Default.GOPATH, build.IgnoreVendor); err != nil {
		return suite.Name
	}
	return fileFileForTest(suite, test, pkg, strip)
}
