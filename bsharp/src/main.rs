use fset::FSet;
use interp::Interp;
use ir::IR;

const SOURCE: &'static str = r#"
# Block comment
# Small comment #
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
      [DEFINE !j [+ !j 1]]
    ]]
    [DEFINE !i [+ !i 1]]
  ]]
  [RETURN !out]
]]

[FUNC @PRINTBOOL [[PARAM !a [BOOL]]] [
  [IF !a [
    [PRINT "true"]
  ] [
    [PRINT "false"]
  ]]
]]

[FUNC @testenum [[PARAM !enum $OPTION_A]] [
  [MATCH !enum 
    [CASE !val $NONE [
      [PRINT "IT'S NONE"]
    ]]
    [CASE !val $A [
      [PRINT "IT'S A"]
      [PRINT [@concat [ARR "VAL 1 is" [GET !val 0]]]]
      [PRINT [@concat [ARR "VAL 2 is" [GET !val 1]]]]
    ]]
  ]
]]

[FUNC @main [] [
  [DEFINE !i 0]
  [WHILE [< !i 10] [
    [DEFINE !i [+ !i 1]]
    [PRINT "HI"]
  ]]

  [@PRINTBOOL [NOT [& [> 1 1] [| [< 1 1] [> 1 1]]]]]
  [DEFINE !enum [NEW $OPTION_A [NEW $NONE]]]
  [DEFINE !v [GET !enum $NONE]]
  [DEFINE !b [NEW $B [: .a "Hello"] [: .b "World"]]]
  [PRINT [GET !b .a]]
  [SET !b [: .a "Hi"]]

  [@testenum !enum]

  [DEFINE !enum [NEW $OPTION_A [NEW $A "Hello" "World"]]]

  [@testenum !enum]

  [MATCH "HELLO"
    [CASE "HELLO" [
      [PRINT "HI"]
    ]]
  ]

  [MATCH 'H' 
    [CASE 'H' [
      [PRINT "HI"]
    ]]
  ]

  [DEFINE !box [BOX "HELLO"]]
  [@PRINTBOOL [PEEK !box $STRING]]
  [PRINT [UNBOX !box $STRING]]

  [PRINT "IS 1+1=2?"]
  [@PRINTBOOL [= [@add 1 1] 2]]
]]
"#;

fn main() {
    // Fset
    let mut fset = FSet::new();
    fset.add_file_source("main.bsp".to_string(), SOURCE.to_string())
        .unwrap();

    // IR
    let mut ir = IR::new(fset);
    ir.build().unwrap();
    if ir.errors.len() > 0 {
        for err in ir.errors {
            println!("{:?}\n", err);
        }
        return;
    }

    // Run
    let mut interp = Interp::new(ir, Box::new(std::io::stdout()));
    let mainind = interp
        .ir
        .funcs
        .iter()
        .position(|f| f.name == "@main")
        .unwrap();
    let res = interp.run_fn(mainind, Vec::new()).unwrap();
    println!("OUTPUT: {:?}", res);
}
