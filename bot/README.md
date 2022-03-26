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

## Buttons
The `BUTTON` function creates a button object, which is just a `MAP{STRING}STRING`. It accepts four arguments:
| Index | Type | Name | Description |
| - | - | - | - |
| `0` | `STRING` | Label | The text on the button. |
| `1` | `STIRNG` | ID | The ID of the button (used later). |
| `2` | `INT` | Color | The Color of the button. See the table [here](#button-colors) for details. |
| `3` | `BOOL` | Disabled | Whether the button is disabled or not. |

The `BUTTONS` function accepts input. It accepts a 2-dimensional array of `BUTTON` objects. Each element in the array is a row of buttons. The return value is the ID of the button pressed. 

The user can be requested to press a button as such:
```scala
[PRINT 
  [BUTTONS
    [ARRAY
      [ARRAY
        [BUTTON "Click Me!" "clicked" 4 FALSE]
      ]
    ]
  ]
]
```
It will print `clicked` (the ID of the button) when the button is clicked.

A grid of buttons can be generated as such:
```scala
[DEFINE btns [MAKE MAP{STRING}STRING{}{}]]
[DEFINE i 0]
[WHILE [COMPARE [VAR i] < 3]
  [DEFINE j 0]
  [DEFINE row [MAKE MAP{STRING}STRING{}]]
  [WHILE [COMPARE [VAR j] < 3]
    # Add btn
    [APPEND [VAR row] [BUTTON "🔲" [CONCAT [STRING [VAR i]] "," [STRING [VAR j]]] 1 FALSE]]

    [DEFINE j [MATH [VAR j] + 1]]
  ]
  [APPEND [VAR btns] [VAR row]]
  [DEFINE i [MATH [VAR i] + 1]]
]

[PRINT 
  [BUTTONS 
    [VAR btns]
  ]
]
```
It will print out the position of the button ont he grid when a button is pressed.


### Button Colors
The following button colors are available:
| ID | Color |
| - | - |
| 1 | Primary |
| 2 | Secondary |
| 3 | Danger |
| 4 | Success | 



## Additional Tags
The following tags are available in the bot in addition to the builtin tags:

| Tag | Description |
| - | - |
| `[USERID]` | Returns the User ID of the user running the tag |
| `[INPUT prompt]` | Prompts the user to provide input, and returns the provided input |
| `[MULTIPLAYER enabled]` | If enabled, allows multiple people to respond to buttons and `INPUT` |