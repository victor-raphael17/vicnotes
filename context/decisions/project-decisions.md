The project will have three main versions, local, cloud and high traffic cloud.

For the local version, I won't use Redis, since it's not needed due to the low traffic.

For the low traffic cloud version, I won't use Redis, since it's not needed and it's supossed to be as cheap as possible.

For the high traffic cloud version, I will use Redis, since it's needed due to the high traffic.

I won't be using node.js for the backend. I chose Go because it has better maintainability.

Determinated the use of cloud services for the cloud versions.
