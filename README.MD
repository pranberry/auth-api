Creates and returns a JWT and allows a user to access protected endpoints with a valid JWT
Endpoints:
- /register
    - register a user
- /login
    - login and retrieve a JWT
- /secret
    - validates a legit JWT and sends the client some guarded data
This project is to learn GO-lang.
=========================================================================================

ToDo: 
- add tests for jwt creating and validating
- fix failing tests [✅]
- improve token validation
- Add a postgres database [✅]
    - this, in general, could be the backend for all JWT Rest projects (any lang) [✅]
- Instead of a simple secret key for signing, impliment public-private infra
- rename "auth" folder to token
- store secret key in db, and get it from db

=========================================================================================
So what did I learn?
- Go is a pretty fun language!
- Opinionated, but not overbearing
- VSCode integrations is resonable
    - linting feels a bit agro...but hey, we get compiled
    - intellisense is quite good...atleast for the main go stuff. some of the imported packages leave a bit to be desired
    - debugger is quite good
- packaging felt quite straight forward
- Interesting to see pointers get used. Haven't used them outside of acedemic contexts before
    - I'm not having to do a lot of extra work to edit objects right in memory
        - it doesn't allow "memory walking" like with c-lang (can't do: &var + 4)
    - with something like java, you'd never be doing in-memory editing
    - I wonder how much less the overhead actually is when compared to some lang like Java?
    - Go vs Java
        - Go does Ahead of time compiling (AOT). compile once, predictable performace, no dynamic tuning.
        - Java does runtime optimization. It has the JIT. 
            - You compile, but overtime the code is observed
            - The performace is fine-tuned.
            - frequently used paths are written in machine code, function inlining, speculative execution, etc
            - Frequently used code get compiled again, optimized and replaced live! holy cow!
        - Go on the other hand has lower memory footprint, faster startup (and faster cold starts), low cpu idle
            - a go-routine uses tiny bit of memory
        - Who knew garbage collector tuning is a thing!?
- reciever vars != input args
    - func (T Type) FuncName(V ValueType) (ReturnType1, ReturnTypeN) {...}
    - the (T Type) are the reciver args. they bind the function to that type
    - (V ValueType) are inputs you can use, but, like with funcs in general, are optional
- also, the whole "var_name var_type" syntax is a bit wacky. why not just do "var_type var_name", like every other lang?
- no overloading ':(', but types switches can be useful
- variadic functions, funcs which accept a variable number of args
    - to clarify, the last arg (and only the last arg) can be declaired to be a variadic arg
    - it says you can pass any number of values for that arg
    - Example: 
        - func Foo(val ...int)
        - this says "func Foo will take any number of int values for val"
    - inside the func, the values are treated as a slice
    - Related: you can "spread" a slice
        - that is, if you have a slice "slc := []int{1,2,3}"
        - you can "spread" it into a function by doing:
            - Foo(slc...)
    - neato!
- that native error type is hella cool
- remembering to include github package names, and remembering to include them can be difficult

- This is actually a great little project for learning a lang.
    - hiting a lot of core functionality of a lang here
    - learned a LOT
    - need to look indepth into the underlying philosophy of Go
        - how is it different from java and python?
        - what problem were its creators trying to solve?

What more to checkout with GO-lang:
- type assertions
- No interhitance, but we've got interfaces...
    - composition vs inheritance?
    - generic interface{}
    - type assert on an interface
    - type switches
- what is a runtime panic? when do they pop up?

Best Practices: what i learned as i went along
(May 01, 2025)
- structs over maps (always)
- defer stmt.Close() when preparing statement
    - very interesting. when you defer it closes the statement when the func ends
    - no need to explictly close elsewhere (unless you need to for performance reasons)
- Queried rows must be scanned
    - don't query a column, that you don't want scanned
    - though you can ignore queried columns
- Best practice for returns is (T, err). not (T, bool)
    - (T, err) is much more specific
    - (T, bool) is okay for in-mem lookups (checking if val exists in map)
- doing the cutover to the db was fairly straight forward...painless even
- retro fitting the tests was really fukin easy. took out a lot of bloat too
    - just needed to add a test harness
- apparently instead of returning a generic 'interface{}', you can return an 'any' type...as of go1.18