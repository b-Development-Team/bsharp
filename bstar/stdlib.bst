# THESE ARE ALREADY B* LIBRARIES, BUT IF THEY ARE REMOVED THIS IS THE SOURCE

# baselib by cloverpepsi
[DEFINE digits 0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/]
[FUNC baseconv [ARRAY num b1 b2] 
[BLOCK 
    [DEFINE in 0]
    [TRY [# [INT [VAR b1]]] [THROW "BASE1 argument is not an integer."]]
    [TRY [# [INT [VAR b2]]] [THROW "BASE2 argument is not an integer."]]
    [IF [COMPARE [VAR b1] > 64] [THROW "BASE1 argument is greater than 64."]]
    [IF [COMPARE [VAR b2] > 64] [THROW "BASE2 argument is greater than 64."]]
    [DEFINE num [STR [VAR num]]]
    [MAP [DEFINE in [+ [VAR in] [* [^ [VAR b1] [VAR map.i]] [IF [COMPARE [INDEXOF [SLICE [VAR digits] 0 [VAR b1]] [VAR num [VAR map.i]]] = -1] [THROW "NUM argument contains unexpected digits."] [INDEXOF [SLICE [VAR digits] 0 [VAR b1]] [VAR num [- -1 [VAR map.i]]]]]]]] [RANGE 0 [LENGTH [VAR num]] 1]]
    [DEFINE out ""]
    [DEFINE p 1]
    [WHILE [COMPARE [VAR in] > 0]
    [BLOCK
        [DEFINE temp [% [VAR in] [^ [VAR b2] [VAR p]]]]
        [DEFINE in [- [VAR in] [VAR temp]]]
        [DEFINE out [CONCAT [INDEX [VAR digits] [INT [/ [VAR temp] [^ [VAR b2] [- [VAR p] 1]]]]] [VAR out]]]
        [DEFINE p [+ [VAR p] 1]]
    ]]  
    [IF [COMPARE [VAR out] = ""] [DEFINE out 0]]

[RETURN [VAR out]]]]


# bsharplib by nv7
[IMPORT baselib]

[DEFINE c "                                 ! #$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"]
[FUNC XOR [ARRAY a b] [BLOCK
  [DEFINE l [STR [baseconv [VAR a] 10 2]]]
  [DEFINE r [STR [baseconv [VAR b] 10 2]]]
  [IF [COMPARE [LENGTH [VAR l]] > [LENGTH [VAR r]]]
    [DEFINE r [CONCAT [REPEAT "0" [MATH [LENGTH [VAR l]] - [LENGTH [VAR r]]]] [VAR r]]]
    [DEFINE l [CONCAT [REPEAT "0" [MATH [LENGTH [VAR r]] - [LENGTH [VAR l]]]] [VAR l]]]
  ]

  [DEFINE out ""]
  [DEFINE i 0]
  [WHILE [COMPARE [VAR i] < [LENGTH [VAR l]]] [BLOCK 
    [DEFINE lv [INT [INDEX [VAR l] [VAR i]]]]
    [DEFINE rv [INT [INDEX [VAR r] [VAR i]]]]
    [DEFINE res [MATH [MATH [VAR lv] - [VAR rv]] * [MATH [VAR lv] - [VAR rv]]]]
    [DEFINE out [CONCAT [VAR out] [VAR res]]]
    [DEFINE i [MATH [VAR i] + 1]]
  ]]

  [RETURN [INT [baseconv [INT [VAR out]] 2 10]]]
]]

[FUNC BOR [ARRAY a b] [BLOCK
  [DEFINE l [STR [baseconv [VAR a] 10 2]]]
  [DEFINE r [STR [baseconv [VAR b] 10 2]]]
  [IF [COMPARE [LENGTH [VAR l]] > [LENGTH [VAR r]]]
    [DEFINE r [CONCAT [REPEAT "0" [MATH [LENGTH [VAR l]] - [LENGTH [VAR r]]]] [VAR r]]]
    [DEFINE l [CONCAT [REPEAT "0" [MATH [LENGTH [VAR r]] - [LENGTH [VAR l]]]] [VAR l]]]
  ]

  [DEFINE out ""]
  [DEFINE i 0]
  [WHILE [COMPARE [VAR i] < [LENGTH [VAR l]]] [BLOCK 
    [DEFINE lv [INT [INDEX [VAR l] [VAR i]]]]
    [DEFINE rv [INT [INDEX [VAR r] [VAR i]]]]
    [DEFINE res [MATH [MATH [VAR lv] + [VAR rv]] - [MATH [VAR lv] * [VAR rv]]]]
    [DEFINE out [CONCAT [VAR out] [VAR res]]]
    [DEFINE i [MATH [VAR i] + 1]]
  ]]

  [RETURN [INT [baseconv [INT [VAR out]] 2 10]]]
]]

[DEFINE h [ARRAY "HEAP"]]
[FUNC HEAPADD [ARRAY v] [BLOCK
  [DEFINE res [LENGTH [VAR h]]]
  [DEFINE h [CONCAT [VAR h] [ARRAY [VAR v]]]]
  [RETURN [VAR res]]
]]
[FUNC HEAPGET [ARRAY ind] [INDEX [VAR h] [VAR ind]]]
[FUNC HEAPSET [ARRAY ind v] [DEFINE h [SETINDEX [VAR h] [VAR ind] [VAR v]]]]

[FUNC STRSETINDEX [ARRAY v ind val] [BLOCK 
  [DEFINE l [SLICE [VAR v] 0 [VAR ind]]]
  [DEFINE r [SLICE [VAR v] [MATH [VAR ind] + 1] [LENGTH [VAR v]]]]
  [RETURN [CONCAT [VAR l] [VAR val] [VAR r]]]
]]