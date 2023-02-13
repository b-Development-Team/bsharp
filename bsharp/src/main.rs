use fset::FSet;
use ir::IR;

const SOURCE: &'static str = r#"
# Block comment
# Small comment #
[TYPE $STRING [ARRAY [CHAR]]]

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
}
