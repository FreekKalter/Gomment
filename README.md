# Gomment

A streamfilter to comment out some common errors when writing *Go* code.

- unused variables
  - simple var declaration, the whole line gets commented out
  - on a multiple var declaration on one line it gets replaced with _
- unused imports
  - unused imports get commented out by line
- no new variables on left side of :=
  - if all variables on left side of := are declared allready it becomes a =


## Example
*this*
    
    package main
  
    import "fmt"
    import "os"
  
    func main(){
        var a int
        var b,c int
        d := 1
        e := 1
        d, e := 2,2
        fmt.Println(d,e)
    }
*becomes*
    
    package main
  
    import "fmt"
    //import "os"
  
    func main(){
        //var a int
        var b,c int
        var d := 1
        var e := 1
        d, e = 2,2
        fmt.Println(d,e)
    }
