##buildconstants 
### for putting build info (and other compile-time values) in Go constants

Sometimes it can be useful to include constants that were available at compile-time, the build number, git branch, etc.
This utility will generate a go file with constants that are the output of various commands.
While you could do something similar with linker flags, this is intended for constants (and doesn't add args to your build command).

The goal is to be able to include this in a go generate comment so the file will be automatically created as necessary.

It does include a constant `MustRunBuildConstants` so that the file with the go generate statement can be used to force a useful compiler error 

##Example of using with go generate

example.go:

    //go:generate buildconstants
    
    // this will generate a compiler error to remind you to run it
    var _ = MustRunBuildConstants 
    


##Example:
  Given a text file with shell commands, like so:
    
    GOVERSION = go version
    BUILD_NUMBER = ${BUILD_NUMBER}
    

  Running `go generate` before building will evaluate the commands and variables in the shell.
    
  This will be used to generate something like:
    
    package currentPackage
    
    const (
        GOVERSION = "go version go1.6.2 linux/amd64"
    )
    
  Which then allows you to use `currentPackage.GOVERSION` in your code.

##Arguments 
   
    -o output file name (defaults to buildconstants_generated.go)
    -package the package the file will be in (defaults to current package)
    -i input file of commands 

##Note:
Don't use this to put date or time related data in constants, your builds should still be [reproducible](https://reproducible-builds.org/), which is necessary for using tools like [bazel](http://bazel.io).

##License 

MIT 
