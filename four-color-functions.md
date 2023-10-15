# What Color is your Go Function; Red, Green, Blue or Cyan?

I recently read this excellent blog post entitled _[Rust Threads vs. GoRoutines](https://shane.ai/posts/threads-and-goroutines/)_ by [Shane Hansen](https://www.linkedin.com/in/shanemhansen/) — which I found linked from the [Golang Weekly](https://golangweekly.com/) that I think ever self-respecting Gopher should be subscribed to — where the blog post discussed creating a million threads in Rust vs. a million Goroutines and the respective overhead required for each.

And then in the post's penultimate paragraph I came across this comment with a reference to [Bob Nystrom](https://twitter.com/munificentbob)'s seminal _[What Color is Your Function?](https://journal.stuffwithstuff.com/2015/02/01/what-color-is-your-function/)_ post:

> _"It’s possible for a call to .Read() to submit a non-blocking I/O request and then cooperatively switch to the next goroutine much like a async Rust function **but without having the colored function problem** that often leads to library bufurication [sic]."_

While it is can be considered true that Go avoid bifurcation related to its Goroutines, reading that reminded me of another aspect of Go which I have agonized over for several years now, an aspect that from my perspective results in Go having *four (4) different function colorizations* albeit they are not quite as mutually-exclusive as Javascript's `async` vs. non-`async` functions.

## Go's Four Different Func Colorizations

Unlike Javascript's two [Red and Blue functions](https://journal.stuffwithstuff.com/2015/02/01/what-color-is-your-function/#:~:text=Every%20function%20has%20a%20color), Go has four (4) different types of functions, which we'll identify with these four (4) colors:

1. Red
2. Blue
3. Green
4. Cyan

Note that #4 in color theory is the combination of #2 and #3; that was intentional.

## Ergonomic, Functional, and Robust: Pick _(at most)_ Two

There are three (3) different aspects affected by the color of your function, and we'll find when mixing colored functions you can _at best_ only get two, but often only one.

### 1. Ergonomic
Simply put, can you call the function in an expression, or do you have to call it as a statement in order to assign multiple return values. This rears its head most frequently when you would like to use an `if` or `for` condition but the function you are calling returns more than one result, e.g.

```go
if Ergonmic() {
    printlin("Easy Peasy!")
}

// vs.

result,err := NotErgonmic() 
if err!=nil {
    return err
}
if result {
    printlin("That was a lot of effort!")
}
```

Note I am not calling out the need to check the error — I strongly believe in the benefits of that aspect of Go — I am only denoting it is not possible to check the error without first having to assign both return values to a variable. I think there should be a better way, and one that does not make Go a significantly more complex language.

### 2. Functional
Here I am specifically referring to **one** specific functionality — the ability to manage called functions by passing in a `context.Context`.  Passing contexts allows for callers to:

1. Request cancellation, 
2. Pass in deadlines and timeouts, as well as
3. Pass in request-scoped data in a commonly-accessible form, e.g. session IDs, auth tokens, etc.

Typically, contexts are passed in as a first parameter, e.g.:   

```go
func ProcessWidgets(ctx context.Context, widgetChan chan Widget) error {
   for {
      select {
      case <-ctx.Done():
         return ctx.Err()
      case widget, ok := <-widgetChan:
         if !ok {
            return nil
         }
         // Do something with the widget
      }
   }
}
```   
Problems can arise when funcs has been written that do not accept a context, and later we realize that they need context functionality, such as this contrived example where the `.Process()` method of the widget does not accept a change, but we find out later it needs to call the time-consuming method `.AssociateItems()` which could benefit from cancellation and/or deadline or timeout functionality.

Notice how `.Process()` has to create a new context, which disconnects it from the top-level context and thus means that any cancellations sent via the outer func will not make their way into the inner func and cancellation will not work until all items for a widget have been associated:

```go
func ProcessWidgets(ctx context.Context, widgetChan chan Widget) error {
   for {
      select {
      case <-ctx.Done():
         return ctx.Err()
      case widget, ok := <-widgetChan:
         if !ok {
            return nil
         }
         err := widget.Process()
         if err != nil {
            return err
         }
      }
   }
}
func (w *Widget) Process() error {
   // process stuff
   ctx := context.Background() // NOT SAME AS DoSomething()'s context
   err := w.AssociateItems(ctx)
   if err != nil {
      return err
   }
   // process more stuff
   return nil
}
func (w *Widget) AssociateItems(ctx context.Context) error {
   for {
      select {
      case <-ctx.Done():
         return ctx.Err()
      case item, ok := <-w.ItemChan:
         if !ok {
            return nil
         }
         // Do something with the item
      }
   }
}
type Widget struct {
   ItemChan chan Item
}
type Item struct {}
```

The _"fix"_ is to change the method signature of `.Process()`, break any existing calling code, and possibly break any interfaces that use that signature, but that fix can ripple downstream and cause an entire host of other problems which in my opinion do not make for good software engineering.

Like with our prior section on ergonomics, I think there should be a better way to handle contexts that do not require the choice of adding a first context parameter to every function and then dealing with the fallout when a context wasn't added, but it was later discovered to be needed.

### 3. Robust
The third and final aspect we'll look at for this essay is robustness in Go, which is inversely related to ergonomics, at least currently.

In Go [errors are values](https://go.dev/blog/errors-are-values), and it is idiomatically considered a best practice to capture  errors as a value and then take action on that value; be it to correct, retry, and/or log the error, as applicable. 

Overwhelmingly the convention for communicating an error value from the called function back to the caller is to return the error value as the last of one or more return values. You can see many examples of funcs returning an error from our prior examples.

The problem for robustness in Go arises when during maintenance you discover you need to call a public function that returns an error value _**after**_ your code's public API has been released for consumption by others. 

Consider if a `Widget` had a `.Price() *Price` method which when initially written was a static value stored in the widget that was set when the widget first instantiated. The method clearly does not generate an error, so obviously there was no need for it to return an error value, and our `.FormattedPrice()` method below calls it as such:

```go
type Widget struct {
   Name string
   price *Price
   file *os.File
}
func (w *Widget) Process() error {
   // process stuff
   _,err := w.file.Write(w.FormattedPrice())
   if err != nil {
      return err
   }
   // process more stuff
   return nil
}
func (w *Widget) FormattedPrice() string {
   // Call w.Price() and fmt.Sprintf() will see it has 
   // a .String() method and call it.
   return fmt.Sprintf("%s: %s",w.Name,w.Price())
}
func (w *Widget) Price() (p *Price) {
   return w.price
}
func (p *Price) String() string {
   return fmt.Sprintf("%s%d%02d",
       p.Currency.Symbol(),p.Major,p.Minor)
}
type Price struct {
   Major int
   Minor int
   Currency CurrencyType
}
type CurrencyType int
const (
   USD CurrencyType = iota
   EURO
)
func (c CurrencyType) Symbol() string {
   switch c {
   case USD:
      return "$"
   case EURO:
      return "€"
   }
   return "?"
}
```

Then imagine your COO decides to pursue a real-time pricing strategy meaning your widget needs to reach out to external systems to calculate a price, and the COO will not even let you cache it.  

NOW your method signature needs to become `.Price() (*Price,error)` and every function that calls your price method either needs to now return an error, or log the problem and retry or correct it as applicable. We now have to change our signatures for our `Widget` methods `.Process()`, `.FormattedPrice()` and `.String()`:

```go
func (w *Widget) Process() error {
   // process stuff
   fp, err := w.FormattedPrice()
   if err != nil {
      return err
   }
   _,err = w.file.Write(fp)
   if err != nil {
      return err
   }
   // process more stuff
   return nil
}
func (w *Widget) FormattedPrice() (string,error) {
   p,err := w.Price()
   if err != nil {
      return "",err
   }
   return fmt.Sprintf("%s: %s", w.Name,p),nil
}
func (w *Widget) Price() (*Price,error) {
   p,err := RealTimePrice(w)
   if err != nil {
      return nil,err
   }
   return p,err
}
func RealTimePrice (w *Widget) (p *Price,err error) {
   // do your realtime pricing here
   return p,err
}
```

Even if everyone who works with your source code is on your own team and/or within the same company, this kind of code thrashing can be disruptive and results in pull-requests that can be very hard to land because of how broadly they can infect a codebase. Not to speak to the problems created for 3rd parties that are currently using published packages.


## The Matrix of Our Discontent

What might be the distinction between these different function colorizations? 

Before I spell it out, let's look at how they relate to each other when calling one from another:

1. Green, Blue and Cyan functions can call Red functions ergonomically and with no loss of functionality or robustness.
2. Green functions can call Cyan functions with no loss of functionality.
3. Blue functions can call Cyan functions with no loss of robustness.
4. Blue functions can call Green functions ergonomically, but they **forfeit robustness**.
5. Green functions can call Blue functions but they **forfeit functionality**.
6. Red functions can call Green, Blue or Cyan functions but **forfeit functionality and/or robustness**.
7. Cyan functions can call any color function with no loss of functionality and/or robustness.
8. Blue and Cyan functions **cannot be used ergonomically** in an expression or assignment.
9. Red functions can be used ergonomically in an expression and in assignment.
10. Changing the color of any function **breaks** its method signatures and any related interfaces.

## What Does Each Color Represent?

The distinctions I am making relate to whether or not a function:

1. Accepts a `context.Context` as its first parameter, and/or
2. Returns an `error` as its last return value.

If a Go function returns an `error` its use can be made robust, but it cannot be called ergonomically.

Here it is in chart form:

|Color| ERGONOMIC<br>Used in expressions<br>and assignments |  FUNCTIONAL<br>Accepts a<br>`context.Context`  |  ROBUST<br>Returns an<br>`error`  |
|---|:---------------------------------------------------:|:----------------------------------------------:|:---------------------------------:|
|Red |                         Yes                         |                       No                       |                No                 |
|Green|                         Yes                         |                      Yes                       |                No                 |
|Blue |                         No                          |                       No                       |                Yes                |
|Cyan |                         No                          |                      Yes                       |                Yes                |

## The Worst Aspect of Function Colorization

Are **non-ergonomic calls** the worst aspect of Go function colorization? While inability to call a function ergonomically is an annoyance, it really does not affect the quality of the software that can be developed. Instead, that is just a nice-to-have that we often don't get to have when programming in Go, especially if we are concerned about robustness.

What about **lack of functionality**? Not really, because adding a `context.Context` parameter to gain the functionality is not too hard to add when needed. It only needs to be added all the way through the call stack between the two points about where it is required.

Is **lack of robustness** the killer here? Again, not really. Just like `context.Context`, adding a returned `error` is not terribly difficult, you again just need to add it all the way through the call stack between the two points about where it is required, and you need to modify any call locations where one of the functions in the stack was used as an expression or in an assignment. Annoying yes, but doable.

So if not non-ergonomic calls, lack of functionality, or lack of robustness what is the worst aspect of function colorization in Go?

_(Of course, the lack of ergonomic calling when errors are returned is an incentive for developers to rationalize an error will likely never occur, so they can ignore it and be able to call the function ergonomically. But I digress.)_

### Broken Function Signatures
The worst aspect of function colorization in Go are how it results in **broken function signatures** — and **especially in interfaces** — after code has been shipped and used by 3rd parties. 

Sure, we can add new function signatures and deprecate the old ones, but that can add considerable cognitive load when using a package and many functions are really hard to find an alternate name for.

It is not quite as bad for _internal-use only code_, but still it can create a large amount of code churn, disruption on a team of developers when changes are made, requiring a lot of rote code review, and making some pull requests really hard to land because of the number of files that can be effected.

All-in-all I find Go to be a programming language that is **hard to future-proof** without adding to every function `context.Context` as a parameter and the return of an `error`. When time comes that a `context.Context` or an `error` are required to be added for functionality and/or robustness to functions that do not already have them, **Go makes it impossible** to do so with existing named funcs and named interfaces and instead requires deprecating the funcs and interfaces and replacing them with new names. 

Or one can just ignore the fallout which is what happens for most internal-only Go projects.     

## Can Go's Functions be Unified into just One Color?

I think so. 

However, rather than discuss how that might be possible or make a language proposal for such a change I have discovered one must make certain others who read and comment on proposals are damn clear on the issues at hand, and that most agree it is a problem. 

Otherwise, any proposal made without first ensuring a shared understanding and agreement will fall on deaf ears in the best case, and the in worse cases will receive hostility from those who do not recognize and/or appreciate the problem one is trying to solve.      

### Can you help?
If you agree that this is a problem and that finding a solution would be helpful, can you comment below but also share it places where Go developers are likely to see it

And if you disagree that it is a problem, I would love to hear your feedback in the comments as well. If so maybe I will come to appreciate your perspective more than I appreciate my current one.