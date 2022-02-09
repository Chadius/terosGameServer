# Why
The underlying package may panic or throw errors. We need to be prepared to catch those issues.

# What is it
Panics call defer, giving us a chance to recover and return the error.
With an error, the server will return a 500 status code.

# What happens next
Clients will be able to better respond to errors.