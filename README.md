##buildconstants - used for putting build info in Go constants

##Example:
  A variable (which are put in a shell environment) and a command which generates the value.
   
    GOVERSION=`go version`
    
  Evaluating might give something like:
    `go version go1.6.2 linux/amd64`

  This will be used to generate 
    
    `package currentPackage
    
    const (
        GOVERSION = "go version go1.6.2 linux/amd64"
    )`
    
  Which then allows you to use `currentPackage.GOVERSION` in your code.

##License 

MIT 
