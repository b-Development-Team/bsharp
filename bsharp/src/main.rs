use fset::FSet;
use ir::IR;

const SOURCE: &'static str = r#"
# Block comment
# Small comment #
[TYPE $STRING [ARRAY [CHAR]]]
[TYPE $A [TUPLE $STRING $STRING]]
[TYPE $B [STRUCT
  [FIELD .a $STRING]
  [FIELD .b $STRING]
]]
[TYPE $C [TUPLE [GENERIC $T $ANY] [G $B $T] $NONE]]
[TYPE $D [ARRAY [GENERIC $T $ANY]]]
[TYPE $ENUM [ENUM 
  $A
  $B
]]
[TYPE $NONE [STRUCT]]

[TYPE $OPTION_A [ENUM $NONE $A]]
[TYPE $ANY [INTERFACE]]
[TYPE $E [STRUCT 
  [GENERIC $T $ANY]
  [FIELD .a $T]
]]

[TYPE $F [INTERFACE
  [GENERIC $T $ANY]
  [INT] # Can be an int
  [FLOAT] # Can be a float
  [TUPLE $STRING [INT]] # Can be a tuple that starts w/ a string and INT
  [ENUM $B $C] # Can be an enum that contains types B and C (but can contain more)
  [STRUCT [FIELD .a $STRING]] # Can be a struct that contains field a
]]


[FUNC @hello [] [
    [PRINT "Hello, World!"]
]]

[FUNC @add [[PARAM !a [INT]] [PARAM !b [INT]] [RETURNS [INT]]] [
    [RETURN [+ !a !b]]
]]

[FUNC @concat [[PARAM !pars [ARRAY $STRING]] [RETURNS $STRING]] [
  [DEFINE !len 0]
  [DEFINE !i 0]
  [WHILE [< !i [LEN !pars]] [
    [DEFINE !len [+ !len [LEN [GET !pars !i]]]]
    [DEFINE !i [+ !i 1]]
  ]]
  [DEFINE !out [NEW $STRING !len]]
  [DEFINE !i 0]
  [WHILE [< !i [LEN !pars]] [
    [DEFINE !j 0]
    [DEFINE !par [GET !pars !i]]
    [WHILE [< !j [LEN !par]] [
      [APPEND !out [GET !par !j]]
    ]]
    [DEFINE !i [+ !i 1]]
  ]]
  [RETURN !out]
]]

[FUNC @slice [[GENERIC $T $ANY] [PARAM !arr [G [ARRAY [GEERIC $T $ANY]] $T]] [PARAM !start [INT]] [PARAM !end [INT]] [RETURNS [G [ARRAY [GENERIC $T $ANY]] $T]]] [
  [DEFINE !out [NEW [G [ARRAY [GENERIC $T $ANY]] $T] [- !end !start]]]
  # the rest of the code
]]
"#;

fn main() {
    let mut fset = FSet::new();
    fset.add_file_source("main.bsp".to_string(), SOURCE.to_string())
        .unwrap();
    let mut ir = IR::new(fset);
    ir.build().unwrap();

    for err in ir.errors {
        println!("{:?}\n", err);
    }
}
