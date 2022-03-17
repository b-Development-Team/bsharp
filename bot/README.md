# B# Bot
[![Official Server](https://img.shields.io/discord/903380812135825459?logo=Discord&logoColor=ffffff&style=for-the-badge)](https://discord.gg/ecyFpQWFE4) 

Using the B# bot, you can run B# code within Discord!

**Official Discord Server**: https://discord.gg/ecyFpQWFE4

## Main Commands
Use `/run` to run code, `/create` to create tags, and `/edit` to edit tags.

Tags can be imported using
```scala
[IMPORT "<tag id>.bsp"]
```

## Other Commands
Images can be added to tags using `/image`, and descriptions can be added using `/description`. View a tag's info and stats using `/info` and a tag's source using `/source`.

## Database
Tags can store data in between uses using the database system. There are 4 commands:
| Usage | Description |
| - | - |
| `[DB SET key value]` | Sets a key in the database to a value |
| `[DB GET key]` | Returns the value of that key in the database |
| `[DB EXISTS key]` | Returns `"true"` if the key exists, `"false"` if it doesn't |
| `[DB USE id]` | Uses another tag's database. :warning: **NOTE**: Other tags' databases are read-only |

## Additional Tags
The following tags are available in the bot in addition to the builtin tags:

| Tag | Description |
| - | - |
| `[USERID]` | Returns the User ID of the user running the tag |
| `[INPUT prompt]` | Prompts the user to provide input, and returns the provided input |