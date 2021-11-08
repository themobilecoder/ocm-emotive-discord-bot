# Emotive

Emotive is a discord bot written in `Go` made for OCM Discord that can tell you the most interesting messages in a channel. These are based on the following criteria:
- Number of reactions
- Message length (TODO)
- Categories (TODO)

## Installation

Using `task`

```bash
task build
```

or using simple bash
```bash
go build -o bin/emotive main.go handlers.go
```

### Note:
Be sure to set your Discord bot's token and your discord guild id in your environment variables
```bash
export DISCORD_TOKEN=ABCD1234TOKEN
export DISCORD_GUILD_ID=ABCD1234GUILDID
```

## Run
```
$ bin/emotive
```

## Usage
 Run this command in a channel, and expect the bot to DM you the results.
```bash
!bestposts <args (TODO) >
```
## TODOs
- Get the top messages by category e.g.
`!bestposts funny`
`!bestposts reactions`
- Define number of posts to retrieve
`!bestposts top 25`

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)