# What is this?
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fartemis-org%2Frestproxy.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fartemis-org%2Frestproxy?ref=badge_shield)

restproxy is a HTTP reverse proxy for the [Discord](https://discordapp.com) API that includes a ratelimiter.

# Where is this used?
We use this utility on the [Artemis](https://artemisbot.io) bot in production.

# How does this work?
restproxy is based on an AMQP RPC architecture. A consumer listens on queue for API requests, performs them and then 
responds through the publisher with the response from the Discord API.

A request object is created, using the URL, the request type, the headers and the content fields from the struct.

The request is then sent to a queue, which the ratelimiter reads from as instructed to by the Discord API, to not trip
the ratelimit and get a 429.

Finally, the response is then proxied back to the client. This means that HTTP requests may well take several seconds,
and therefore clients using this utility must making all requests asynchronously.

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fartemis-org%2Frestproxy.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fartemis-org%2Frestproxy?ref=badge_large)