# B# Docs
## Table of Contents
- [Introduction](#introduction)
- [Hello World](#hello-world)
- [Data Types](#data-types)
  - [Numbers](#numbers)
  - [Strings](#strings)
  - [Arrays](#arrays)
  - [Maps](#maps)
  - [Type Conversions](#type-conversions)
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
Arrays represent a list of items. Their types are made by putting a `{}` after the type of the elements in the array. For example, an array of integers would be written as `INT{}`.

There are two ways of defining arrays. An array literal can be defined using the `ARRAY` function.
```scala
[ARRAY 1 2 3]
```
An empty array can be initialized using the `MAKE` function.
```scala
[MAKE INT{}]
```
Arrays are pointers. This means that a change to an array, like appending to it, appends to all of its uses. For example:
```scala
[DEFINE a [ARRAY 1 2 3]]
[DEFINE b [VAR a]]
[APPEND [VAR a] 4]
```
In this case, `b` will also have `4` as the last element.

### Maps
Maps represent key-value data. For example, take the following table.

| Key | Value |
| - | - |
| Name | Joe |
| Favorite Color | Blue |
| Favorite Language | B# |

Map types are defined by putting `MAP{`, the key type, a `}`, and the value type. For example, `MAP{STRING}INT` would be the type for a map with a key type of string and a value type of `INT`. Maps can also be nested, for example `MAP{STRING}MAP{STRING}INT` would be a map of a map.

To initialize a map use the `MAKE` function. The `[SET]` and `[GET]` functions can be used to read and write to the map. The `[EXISTS]` function returns a boolean of whether a key exists.
```scala
[DEFINE a [MAKE MAP{STRING}STRING]]
[SET [VAR a] "Name" "Joe"]
[SET [VAR a] "Favorite Color" "Blue"]
[SET [VAR a] "Favorite Language" "B#"]
[PRINT [GET [VAR a] "Favorite Language"]]
```

### Type Conversions
`INT`s, `FLOAT`s, and `STRING`s can all be converted to each other using type conversions. Use the `INT` tag to convert a float or string to an integer, the `FLOAT` tag to convert the other two to a float, and the `STRING` tag to convert an integer or float to a string.

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

## Imports
Use the `IMPORT` function to import other B# files. This adds all functions in that file to the global scope and will run the code in that file. It will also add all global variables in that file to the current scope.

The files that are available to import are implementation-specific.