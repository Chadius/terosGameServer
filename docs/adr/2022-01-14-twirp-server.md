# Why
- A server should handle passing messages between services
- Some services are better written in other languages and environments. Let `go lang` handle network calls and game logic. 

# What is it
A server that will receive YAML/JSON streams. It will then call the terosGameRules package to apply the logic.
Then it will return the response to the client.

# Caveats
The format is very dumb right now; it's just YAML/JSON files.
How do I test this? What do I want to test?
- Twirp is a 3rd party app that writes the logic to manage HTTP requests
- Need to test that the server can access the terosGameRules package.