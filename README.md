# What is this?
restproxy is a HTTP reverse proxy for the [Discord](https://discordapp.com) API that includes a ratelimiter.

# Where is this used?
We use this utility on the [Artemis](https://artemisbot.io) bot in production.

# How does this work?
restproxy listens for all requests, via all HTTP request types, such as GET, POST and PATCH.

restproxy then takes the URL of the request and replaces our domain with discordapp.com, to get the real API endpoint 
on Discord's side.

A request object is then created, using the URL, the same request type, the same headers and the same content.

The request is then sent to a queue, which the ratelimiter reads from as instructed to by the Discord API, to not trip
the ratelimit and get a 429.

Finally, the response is then proxied back to the client. This means that HTTP requests may well take several seconds,
and therefore clients using this utility must making all requests asynchronously.