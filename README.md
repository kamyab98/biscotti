# Biscotti

Fast & Scalable Cookie Matching Service

## How to Use

Biscotti is DB-less and use log file to store matched cookies. You can use Fluentd to produce matched cookie in kafka,
or store it in database.

#### Envs

|Env|Description|
|----|-------|
|PIXEL_REDIRECT_KEY|Query param key of your pixel redirection.|
|COOKIE_KEY|Your user cookie key that you want cookie matching with it.|
|SCHEMA_REGISTRY_URL|Avro Schema registry Address.|
|NETWORK_TO_URLS|Map of network and match url of them.|
|COOKIE_DOMAIN|Biscotti set cookie for users who haven't your token too. This is cookie domain for your cookie.|
|KAFKA_TOPIC|Name of Avro Schema Registry Subject.|

## Concepts:

**Network:** The organization that wants to integrate his users with *You* (advertiser in most scenarios).

**Match Tag:** The tag you must place in a user's browser for the network-initiated Cookie Matching workflow (for
example you can use a cookie as the match tag).
**Network ID:** A string ID uniquely identifying a network account for Cookie Matching and other related operations.
**Your User ID:** A string ID uniquely identifying *Your* users.
**Network User ID:** A string ID uniquely identifying network users.

### Scenario 1: Network Initiated - Network Places *Your* Match Tag

In order to initiate this flow, the bidder must place their match tag such that it renders in the user's browser. A
match tag that only returns the *Your* User ID to the bidder may be structured as follows:

```html
<img src="https://biscotti.you.com/pixel?id=NETWORK_ID"/>
```

then *You* redirects to the network URL with *Your* User ID in query params.

https://<NETWORK_URL>/?user_id=<YOUR_USER_ID>

After the redirect, Network should save the match between *Your* User ID and Network User ID.
Now that the network has the matching table, he knows what is his/her user’s *You* User ID and can integrate with
*Yours* other APIs with *You* User ID.

### Scenario 2: *You* Initiated - *You* Place Network Match Tag

The Pixel Matching tag placed by *You* combines the Network’s Cookie Matching URL with additional parameters the
network can use to populate their match table. For a Cookie Matching URL specified as https://ad.network.com/pixel, it
is structured as follows:

```html
<img src="https://ad.network.com/pixel?user_id=<YOU_USER_ID>"/>
```

The network must respond with redirect to *Your* Cookie Matching Service URL:
https://biscotti.you.com/match/?id=NETWORK_ID&user_id=<USER_ID>

After the redirect, *You* saves the match between *Your* User ID and Network User ID.
Now that *You* have the matching table, s/he knows what the network’s user *Your* User ID so networks can integrate
with their own User IDs.

You can read Google documentation about cookie matching and real time
bidding [here](https://developers.google.com/authorized-buyers/rtb/cookie-guide).







