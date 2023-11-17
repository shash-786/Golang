package main
import "fmt"

//NOTE WALRUS OPERATOR CANNOT BE USED OUTSIDE A METHOD
// num := 42  <---- ERROR
// here var num int = 42

const pi float32 = 3.141592653589236278638763

func main() {
    var username string = "shashank"
    fmt.Println("Name is:",username)
    fmt.Printf("Username if of type:%T\n\n",username)

    var isLogged bool = true 
    fmt.Println("Is the User Logged:",isLogged)
    fmt.Printf("isLogged is of type:%T\n\n",isLogged)

    var smallint uint8 = 255// will overflow for 256 
    fmt.Println("Smallint is :",smallint)
    fmt.Printf("samllint is of type:%T\n\n",smallint)

    
    var smallfloat32 float32 = 254.2132434323// will overflow for 256 
    fmt.Println("smallfloat32 is :",smallfloat32)
    fmt.Printf("smallfloat32 is of type:%T\n\n",smallfloat32)
    
    
    var bigfloat float64 = 254.2132434323
    fmt.Println("bigfloat is :",bigfloat)
    fmt.Printf("bigfloat is of type:%T\n\n",bigfloat)

    //WALRUS OPERATOR
    num := 32.27
    fmt.Printf("num is %f and it's type is %T\n\n",num,num)

    str := "test string"
    //  str = 1 
    //./main.go:37:11: cannot use 1 (untyped int constant) 
    // as string value in assignment
    fmt.Println("Str is -->",str)

    fmt.Println("\nConst Pi is -->",pi)
}




