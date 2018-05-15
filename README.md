job-status-listener
-------------------

Listens over HTTP for job status updates, then publishes them to AMQP.

This service accepts a JSON request body at the `POST {uuid}/status` endpoint which describes a status update for the job with an external ID of `{uuid}`. The body should have three keys, all strings: hostname, message, and state, where state should be one of: submitted, running, completed, or failed. More documentation for the public-facing part of this is available at https://cyverse-de.github.io/misc-api/ .
