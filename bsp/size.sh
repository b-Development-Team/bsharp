# Analyze size of binary using bloaty
clang -g examples/a.c -o a.out
bloaty a.out -d symbols --debug-file=a.out.dSYM/Contents/Resources/DWARF/a.out -n 10000 --domain=file > a.txt
rm -rf a.out.dSYM