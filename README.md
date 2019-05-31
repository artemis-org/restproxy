# What is this?
restproxy is a HTTP reverse proxy for the [Discord](https://discordapp.com) API that includes a ratelimiter.

# Where is this used?
We use this utility on the [Artemis](https://artemisbot.io) bot in production.

# How does this work?
restproxy listens on a Redis list, from which it BLPOPs requests, encoded in JSON, which are unmarshalled to a struct
containing the key information about the request.

A request object is then created, using the URL, the request type, the headers and the content fields from the struct.

The request is then sent to a queue, which the ratelimiter reads from as instructed to by the Discord API, to not trip
the ratelimit and get a 429.

Finally, the response is then proxied back to the client. This means that HTTP requests may well take several seconds,
and therefore clients using this utility must making all requests asynchronously.