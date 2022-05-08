bsharp build bsp/main.bsp bsp/parser.bsp bsp/tokens.bsp bsp/types.bsp bsp/ir/ir.bsp bsp/ir/scope.bsp bsp/ir/build.bsp bsp/ir/stmts.bsp bsp/ir/basic.bsp bsp/ir/ops.bsp -o examples/a.c
clang examples/a.c -o a.out
./a.out