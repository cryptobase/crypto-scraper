# Crypto-scraper #
Scrape public crypto exchange trade information.

# Project layout #

```
/app                Main application
/exchanges          Exchange specific api implementations, using the generic REST client
    /bitfinex
/model              Generic trade data model
/restclient         Generic REST client
```

# Basic approach #

```
Load latest record from disk
If the latest record exists
    Fetch trade data newer than latest record from public exchange api into exchange specific model
Else
    Fetch x recent trade records from public exchange api into exchange specific model
Convert specific model to generic model
Persist new data to disk
```