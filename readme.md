# Retrier is a simple retry-ing mechanism for Golang
Simple, ready to use, suitable for writing / saving something in the background, fire and forget (Haven't developed for returning value though).

## How to use
1. Create an object that implement **Exec()** function that **returns an error** if something is wrong
2. This error will be an indicator whether **Exec()** will be retried or not
3. Create new retrier object, start the "auto retry" by calling **Start()** method

## Example
```golang
import (
    ...
    "github.com/adityarama1210/retrier"
    ...
)

// example struct
type UserActivitySaver struct {
   UserEmail string
   Activity string
}

// example implementation, should be based on you need to do
func (u *UserActivitySaver) Exec() error {
    err := SaveToDB(u.UserEmail, u.Activity)
    if err != nil {
        return err
    }
    
    return nil
}

func main() {
    ...
    // example object / flow
    user := GetUserByEmail("john@gmail.com")
    uas := &UserActivitySaver{
        UserEmail: "john@gmail.com",
        Activity: "Login Success",
    }
    // example object or flow
    
    /* Part to initiate and start the auto retrier */

    // create new retrier for UserActivitySaver object, with maximum retry of 10 times, and using custom logging
    // logging is an optional, put nil if you dont need to
    r := retrier.New(uas, 10, func(err error){
        log.Println("Error on my method", err)
    })
    r.Start() 
    // Start() will trigger the UserActivitySaver Exec method to be executed in the background via go routine, if an error is happening, it will retry up to the max attempt.
    ...
}
```

### How is it running
the retrier will execute Exec() method of the object (in the background via go routine), if no error then it will finish. Otherwise (error) it will try again (with incremented attempt counter) until the attempt is more than or equal to max attempt.