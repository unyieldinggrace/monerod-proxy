# Monerod-Proxy

Monerod-Proxy is a small HTTP proxy server intended to sit in front of a real Monerod node. Your wallet (or other app that needs a node) can talk to the proxy just as if it was talking to the node directly, but with some nice extra features.

### Current Features
* __Automatic Fail-Over__ - The proxy has a config file containing a list of nodes. If the first node in the list stops responding, the proxy will detect this and start forwarding requests to the next node in the list, and so on. If you are using a public node to power your application (or some other node that may have outages), using this proxy layer to automatically switch nodes lets your application keep running without intervention even when a particular node goes down for unplanned/unannounced maintenance, etc.
* __Graceful Redirect__ - If you need to restart your node (e.g. when upgrading to a newer version of the `monerod` software), you can manually tell `monerod-proxy` to disable that node and start redirecting traffic to one of its fallback nodes. When you are done with your maintenance, you just tell `monerod-proxy` to re-enable the node and it will start sending traffic there again. Any clients connected via `monerod-proxy` will see no service interruption. Enable/disable is done via a HTTP API.
* __Debug Logging__ - if you are running a public node, you can put monerod-proxy in front of it and set the log-level to "Debug", then you can watch the requests that are coming in with timestamps and IP addresses to get a better understanding of how your node is being used.

### Potential Future Features
These are some useful features that could be added in the future. If you need these features, monerod-proxy is a good starting point to build upon. These are features that are likely to be very useful for node operators, but they would also add extra complexity if we tried to build them directly into `monerod`. Building these types of convenience features in a proxy layer lets `monerod` stay focused on its core mission: validating the Monero ledger.

* Add TLS support for encrypted traffic. Monerod supports this natively, but monerod-proxy currently doesn't support TLS. Shouldn't be hard to add though, it looks like there is built-in support in Echo, which is the underlying HTTP framework.
* Cache data when response is deterministic from request (e.g. request block by hash). Takes some load off backend nodes, especially if they are public nodes.
* Add a websocket interface to let clients (e.g. cash register app) listen for txns arriving in the mempool.
* Add support for rate limiting.
* Add support for tiered access, so a user could pay an amount of XMR to remove rate limiting for a given access key. This allows public node operators to easily monetize node access.
* Add support for more complex load-balancing between available nodes (instead of just always forwarding request to highest-priority node that is available).

### Getting Started

Clone the repo and run `./build.sh` to build the application inside the `bin` directory. You can run it simply by running `bin/monerod-proxy`. It will load config from the `config.ini` file, so edit that file before running the applcation.

To set the admin password, you can first generate a bcrypt hash by doing a POST request to the proxy:

```
curl http://localhost:18081/proxy/api/generatepasswordhash -d '{"Password":"mysupersecurepassword"}' -H 'Content-Type: application/json'
```
Put the returned hash in the config.ini file. Then you can use your new password to run other administrator requests, such as checking the status:
```
curl http://localhost:18081/proxy/api/status -d '{"Password":"mysupersecurepassword"}' -H 'Content-Type: application/json'
```
Or disabling/re-enabling a node for maintenance:
```
curl http://localhost:18081/proxy/api/setnodeenabled -d '{"Password":"mysupersecurepassword","NodeURL":"http://mynode.com:18081/","Enabled":false}' -H 'Content-Type: application/json'
```

Remember that monerod-proxy doesn't support HTTPS yet, so best to only run these commands from localhost.

## Motivation

Let's say you want to accept XMR for payments on your website. The most obvious way to do this is to run an instance of `monero-wallet-rpc` with a view-only wallet which your web application can use to generate payment addresses and watch for incoming payments.

The problem is that `monero-wallet-rpc` needs to connect to a full Monero node (an instance of `monerod`). Running a full Monero node requires a significant amount of disk space and bandwidth. If you are only expecting to do a handful of XMR sales, then the cost of running a full node on a VPS may be more than you expect to earn by accepting XMR in the first place.

To avoid this cost, you could connect to a public instance of `monerod`, but then you don't control it and it may go down without notice, breaking your payment system. You could also do something creative, like run an XMR node on a home PC and connect your `monero-wallet-rpc` instance to the home PC via a reverse SSH tunnel. But then your payment system will break when you reboot your home PC. In short, if you don't run your own instance of `monerod` in the cloud, then you risk some reliability issues.

This is where `monerod-proxy` comes in. You can run an instance of `monerod-proxy` on your cheap VPS and just connect `monero-wallet-rpc` to that. The proxy will take care of forwarding requests to a real node (like a public node, or your home PC), and will also gracefully handle cases where the real node goes offline (by automatically re-routing requests to another node that is still available). This makes your payment system much more reliable without forcing you to run your own instance of `monerod` on a bigger, more expensive VPS.

Hopefully this enables more people to start accepting XMR for payments with lower overhead costs to get started.

### The Future ###
This basic version is useful for individual merchants who are just looking to get cheap, reliable access to data from monerod to power their XMR checkout process.

But there is potential to build more features into the proxy to power more applications. In particular, I would like to see tiered rate limiting, so that node operators could have users pay for high-speed access to their node. This would allow professional node operators to easily monetize access to their node. This will become more important as the Monero blockchain grows and it becomes slightly less practical/affordable for hobbyists to run their own full node in the cloud.

I would also like to see the proxy provide a WebSocket interface that could be used by a point-of-sale app on a mobile device to get notified of transactions arriving in the mempool. The mobile device could then scan the transactions as they arrive and see if they match the view-key of the merchant. This would make it easy to build a point-of-sale app that just needs a view-key and can do basically-instant scanning of the mempool for a payment that was just requested, making the checkout process much smoother.

But all of that is down the road. I wanted to learn golang, and I wanted a proxy to automatically re-route traffic to another monerod instance if the one I was using went down. I've solved those problems by making this small app.

### Tips ###
If this app is useful to you, feel free to leave a tip.
XMR Tip Jar: 87U3RX6JmHiJfWxGy1nRi1EChpuWXowCqChTAUDh8j6hc9Eg928JreVgJ8DEtW8C97W3MYuh8hKzoHxkpSY4CS7f8NPSHnK
