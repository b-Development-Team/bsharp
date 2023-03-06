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
[TYPE $ENUM [ENUM 
  $A
  $B
]]
[TYPE $NONE [STRUCT]]

[TYPE $OPTION_A [ENUM $NONE $A]]

[FUNC @hello [] [
    [PRINT "Hello, World!"]
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
