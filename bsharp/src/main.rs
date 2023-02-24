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
[TYPE $B [STRUCT 
  [GENERIC $T $ANY]
  [FIELD .a $T]
]]

[TYPE $A [INTERFACE
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
