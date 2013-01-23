# Gomment *(as in Go comment)*

A streamfilter to comment out some common errors when writing *Go* code.

It basicly adresses all the little annoyances expressed by Aldo Cortesi in his
[blog post](http://corte.si/posts/code/go/go-rant.html).
You use it like `gofmt` to format your code, so it takes stdin and spits out 
cleaned on stdout. 

You could integrate it in your editor, I use vim and already had `gofmt` hooked up to clean
up my code with 1 keystroke. Now it first runs through `gomment` and then through `gofmt`.
You could even hook it up to `:w`, to save that last keystroke but thats a personal taste.

## What it does

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
