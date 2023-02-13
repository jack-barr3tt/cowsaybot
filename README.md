# cowsaybot
A silly bot written in Golang to bring [cowsay](https://en.wikipedia.org/wiki/Cowsay) to your discord server.

## Usage
1. Make sure you have `cowsay` installed on your system
2. Create a `.env` file with the following contents:
```js
TOKEN=your_bot_token
```
3. Run `go get ./...` to install dependencies.
4. Run `go run main.go` to start the bot.

## Features
- `moo <text>`: Make a cow say something.
- `moo ping`: Get the current message latency in milliseconds
- `mooquote`: Get a random quote from the [Zenquotes API](https://zenquotes.io/)