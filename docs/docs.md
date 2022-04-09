# B# Docs
## Table of Contents
- [Introduction](#introduction)
- [Hello World](#hello-world)
- [Data Types](#data-types)
  - [Numbers](#numbers)
  - [Booleans](#booleans)
  - [Strings](#strings)
  - [Arrays](#arrays)
  - [Maps](#maps)
  - [Structs](#structs)
  - [Type Definitions](#type-definitions)
  - [Type Conversions](#type-conversions)
  - [Any](#any)
- [Statements](#statements)
  - [If](#if)
  - [While](#while)
  - [Switch](#switch)
- [Variables](#variables)
- [Basics](#basics)
  - [Math](#math)
  - [String Manipulation](#string-maniplation)
- [Functions](#functions)
- [Imports](#imports)
- [Standard Library](#standard-library)
  - [Math](#math)
- [Other Functions](#other-functions)

## Introduction
There are 3 major B.x languages:
- B++
- B*
- B#

B++ is the original version, and was very primitive, with no multiline if statements and loops.

B* is to B++ as C++ is to C. All B++ code is valid B* code, but there are additional features, like loops, multiline if statements, try/catch, and many more functions.

B# is not compatible with B++, but introduces features like type-checking, scoping, and more structural changes. These are the B# docs.

## Hello World
In B#, everything is a tag. Tags are denoted with brackets. The syntax is 
```scala
[NAME arg1 arg2 ...]
```

To print "Hello, World!" in B#, simply do
```scala
[PRINT "Hello, World!"]
```

Comments can be done with `#`. 
```python
# This is a comment
```
There can be inline comments by surrounding code with `#`.

## Data Types
There are 7 data types in B#:

| Type | Description |
| - | - |
| `INT` | Represents an integer. |
| `FLOAT` | Represents a floating-point number (decimal). |
| `BOOL`  | Represents a value of either true or false.
| `STRING` | Represents textual data. |
| `ARRAY` | Represents multiple values in a list. |
| `MAP` | Represents key-value data. |
| `FUNCTION` | Represents a pointer to a function. |
| `STRUCT` | Represents key-value data but with a fixed set of keys. |
| `ANY` | Represents a value of any type. |

### Numbers
There are two types of numbers: integers, and floats. Integers represent whole numbers, whereas floats represent decimals. They can both be negative.
```scala
[PRINT [STRING 1234]]
[PRINT [STRING 1.234]]
```
The `STRING` function can be used to convert an integer or a float to a string, which can be printed.

### Booleans
A boolean can be one of two values: `true`, or `false`. The `COMPARE` function can be used to get booleans. For example, 
```scala
[COMPARE 1 == 0]
```
would be false.

`TRUE` and `FALSE` are constants.

The compare operators are:
| Operator | Description |
| - | - |
| `==` | Tests whether two values are equal. |
| `!=` | Tests whether two values are not equal. |
| `>` | Tests whether the left value is greater than the right value. |
| `<` | Tests whether the left value is less than the right value.
| `>=` | Tests whether the left value is greater than or equal to the right value. |
| `<=` | Tests whether the left value is less than or equal to the right value. |

Logical operations can also be performed on booleans.
```bsp
[AND <a> <b>] # Returns whether both <a> and <b> are true
[OR <a> <b>] # Returns whether <a> or <b> is true
[NOT <a>] # If <a> is true, it returns false, if <a> is false, it returns true
```

### Strings
Strings represent textual data. They are surrounded by quotation marks.
```scala
[PRINT "Hello, World!"]
```

### Arrays
Arrays represent a list of items. Their types are made by putting a `ARRAY{}` after the type of the elements in the curly braces. For example, an array of integers would be written as `ARRAY{INT}`.

There are two ways of defining arrays. An array literal can be defined using the `ARRAY` function.
```scala
[ARRAY 1 2 3]
```
An empty array can be initialized using the `MAKE` function.
```scala
[MAKE ARRAY{INT}]
```
Arrays are pointers. This means that a change to an array, like appending to it, appends to all of its uses. For example:
```scala
[DEFINE a [ARRAY 1 2 3]]
[DEFINE b [VAR a]]
[APPEND [VAR a] 4]
[SET [VAR a] 0 0]
```
In this case, `b` will also have `4` as the last element, and `0` as the first element.

### Maps
Maps represent key-value data. For example, take the following table.

| Key | Value |
| - | - |
| Name | Joe |
| Favorite Color | Blue |
| Favorite Language | B# |

Map types are defined by putting `MAP{`, the key type, a `,`, the value type, and then a `}`. For example, `MAP{STRING,INT}` would be the type for a map with a key type of string and a value type of `INT`. Maps can also be nested, for example `MAP{STRING,MAP{STRING,INT}}` would be a map of a map.

To initialize a map use the `MAKE` function. The `[SET]` and `[GET]` functions can be used to read and write to the map. The `[EXISTS]` function returns a boolean of whether a key exists.
```scala
[DEFINE a [MAKE MAP{STRING,STRING}]]
[SET [VAR a] "Name" "Joe"]
[SET [VAR a] "Favorite Color" "Blue"]
[SET [VAR a] "Favorite Language" "B#"]
[PRINT [GET [VAR a] "Favorite Language"]]
```

### Structs
Structs are like maps, but they only have a fixed number of keys. One advantage however, is that they can store values of multiple types. 

In addition, structs are generally safer because a field is guarenteed to exist, and errors related to accessing nonexistent fields can be caught at compile time. 

The above example could be remade to use structs as follows:

```scala
[DEFINE a [MAKE STRUCT{name:STRING,favoriteColor:STRING,favoriteLanguage:STRING,age:INT}]]
[SET [VAR a] name "Joe"]
[SET [VAR a] favoriteColor "Blue"]
[SET [VAR a] favoriteLanguage "B#"]
[SET [VAR a] age -1]
[PRINT [GET [VAR a] favoriteLanguage]]
```

### Type Definitions
Types can be given name using the `TYPEDEF` declaration. For example:
```scala
[TYPEDEF a INT]
[TYPEDEF arr ARRAY{a}]

[DEFINE val [MAKE arr]]
```

### Type Conversions
`INT`s, `FLOAT`s, and `STRING`s can all be converted to each other using type conversions. Use the `INT` tag to convert a float or string to an integer, the `FLOAT` tag to convert the other two to a float, and the `STRING` tag to convert an integer or float to a string.

### Any
The `ANY` type allows for passing values of multiple types around. It can be checked to be of a certain type using the `CANCAST` tag, and cast using the `CAST` tag. The `ANY` tag can be used to cast a value to type `ANY`. For example, to add 2 integers or floats you could do:
```scala
[FUNC ADD [PARAM a ANY] [PARAM b ANY] [RETURNS ANY]
  [IF [CANCAST [VAR a] INT]
    [DEFINE aval [CAST [VAR a] INT]]
    [DEFINE bval [CAST [VAR b] INT]]
    [RETURN [ANY [MATH [VAR aval] + [VAR bval]]]]
  ]

  [IF [CANCAST [VAR a] FLOAT]
    [DEFINE aval [CAST [VAR a] FLOAT]]
    [DEFINE bval [CAST [VAR b] FLOAT]]
    [RETURN [ANY [MATH [VAR aval] + [VAR bval]]]]
  ]

  [RETURN 0]
]

[DEFINE vala [ANY 1]]
[DEFINE valb [ANY 2]]
[DEFINE res [ADD [VAR vala] [VAR valb]]]
[PRINT [STRING [CAST [VAR res] INT]]]
```
**NOTE**: This will cause a crash in the program if the values are anything other than two `FLOAT`s or two `INT`s!

## Statements
### If
If statements execute code conditionally. For example:
```bsp
[IF [COMPARE 7 == 3]
  [PRINT "7 is equal to 3."] # This executes when the condition is true
ELSE
  [PRINT "7 is not equal to 3."] # This executes otherwise.
]
```
An if statement can also be done without the `ELSE`.
```bsp
[IF [COMPARE 7 > 3]
  [PRINT "7 is greater than 3."]
]
```

### While
While statements repeat until a condition is false. For example, to print the numbers from 1 to 10, a while statement could be used:
```scala
[DEFINE i 1]
[WHILE [COMPARE [VAR i] <= 10]
  [PRINT [STRING [VAR i]]]
  [DEFINE i [MATH [VAR i + 1]]]
]
```

### Switch
A switch-case can be used instead of chaining if statements, when a value wants to be tested. For example, the following could be done with if statements:
```scala
[DEFINE val "Hello"]
[IF [COMPARE [VAR val] == "Foo"]
  [PRINT "The value is foo."]
ELSE
  [IF [COMPARE [VAR val] == "Bar"]
    [PRINT "The value is bar."]
  ELSE
    [IF [COMPARE [VAR val] == "Hello"]
      [PRINT "The value is hello!"]
    ELSE
      [PRINT "What is the value?"]
    ]
  ]
]
```
However, this is messy and inneficient. Instead, it could be done with a switch statement:
```scala
[DEFINE val "Hello"]
[SWITCH [VAR val]
  [CASE "Foo"
    [PRINT "The value is foo."]
  ]

  [CASE "Bar"
    [PRINT "The value is bar."]
  ]

  [CASE "Hello"
    [PRINT "The value is hello."]
  ]

  [DEFAULT
    [PRINT "What is the value?"]
  ]
]
```
These both achieve the same result, however the second is much easier to read and is more efficient.

**NOTE:** The `DEFAULT` statement is optional.

## Variables
Variables are defined and set using the `DEFINE` tag.
```scala
[DEFINE a 1]
```
Would set `a` to 1.

The value of a variable is accessed using the `VAR` tag.
```scala
[DEFINE hello "Hello!"]
[PRINT [VAR hello]]
```
Would print the variable `hello`, after defining it to `"Hello!"`.

Variables can only store values of one type. In the previous example, the variable `hello` would have a type of `STRING`. Due to this, the code below is invalid.
```scala
[DEFINE hello "Hello!"]
[DEFINE hello 1]
```

Variables cannot be accessed outside the scope they are defined in. This means the code below is invalid code, since the variable `a` cannot be accessed outside of the `IF`.
```scala
[IF [COMPARE 1 = 1]
  [DEFINE a "a"]
]
[PRINT [VAR a]]
```

Variables can be re-declared with different types if a new scope is declared. That is why the code below is valid code.
```scala
[DEFINE hello "Hello!"]
[IF [COMPARE 1 = 1]
  [DEFINE hello 1]
  [PRINT [STRING [VAR hello]]]
]
[PRINT [VAR hello]]
```

## Basics
### Math
Use the `MATH` tag to perform math operations. 
```scala
[MATH 1 + 1]
```
Would return `2`.

The supported operators are listed below.

| Operator | Effect |
| - | - |
| `+` | Addition |
| `-` | Subtraction |
| `/` | Division |
| `*` | Multiplication |
| `^` | Exponentiation |
| `%` | Modulo |

### String Maniplation
Use the `CONCAT` tag to concatenate strings.
```scala
[PRINT [CONCAT "Hello, " "World!"]]
```

## Functions
Functions are defined with the `FUNC` tag. Provide paremeters with the `PARAM` tag, and optionally, specify a return type with the `RETURNS` tag. Return values with the `RETURN` tag.
```scala
[FUNC ADD [PARAM a INT] [PARAM b INT] [RETURNS INT]
  [DEFINE result [MATH [VAR a] + [VAR b]]]
  [RETURN [VAR result]]
]
```
Above is an example of a function that takes and returns values. However, functions don't need to do this. The following is also valid.
```scala
[FUNC PRINTHELLO
  [PRINT "Hello"]
]
```

### Array and String Operatiosn
An index of an array or character of a string can be gotten using the `INDEX` tag.
For example
```scala
[INDEX [ARRAY 1 2 3] 0]
```
Would return `1`, the first index of `[1, 2, 3]`. This can also be used to get characters of strings.

Arrays and strings can be sliced using the `SLICE` tag. 
```scala
[SLICE "abcde" 0 3]
```
Would return `abc`. Since arrays are pointers, `SLICE` doesn't return an array. Instead, use
```scala
[DEFINE arr [ARRAY 1 2 3 4 5]]
[SLICE [VAR arr] 0 3]
```
Now, `arr` would be `[1, 2, 3]`.

## Imports
Use the `IMPORT` function to import other B# files. This adds all functions in that file to the global scope and will run the code in that file. It will also add all global variables in that file to the current scope.

The files that are available to import are implementation-specific, except for the standard library files.

## Standard Library
### Math
Import with
```scala
[IMPORT "math.bsp"]
```
Functions:
| Usage | Description |
| - | - |
| `[FLOOR val]` | Returns the floor of a float as an integer |
| `[CEIL val]` | Returns the ceil of a float as an integer |
| `[ROUND val]` | Rounds a float to an integer |
| `[RANDINT lower upper]` | Generates a pseudorandom number between two integers |
| `[RANDOM lower upper]` | Generates a pseuodorandom number between two floats |

### Strings
Import with
```scala
[IMPORT "strings.bsp"]
```
Functions:
| Usage | Description |
| - | - |
| `[JOIN vals joiner]` | Joins vals (an array of strings) with the joiner in between |
| `[SPLIT str sep]` | Splits the string into an array of strings using `sep` as a seperator |

## Other Functions
| Name | Description |
| - | - |
| `[TIME mode]` | Returns the current unix timestamp. **NOTE**: The mode is optional. If provided, it must be `SECONDS`, `MICRO`, `MILLI`, or `NANO` |