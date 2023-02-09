use fset::FSet;
use ir::IR;

const SOURCE: &'static str = r#"
# Block comment
# Small comment #
[PRINT "Hello, World!"]
[+ 1 [- 2 1] [% 1 2]]

[TYPE $STRING [ARRAY [CHAR]]]

# Block
[[PRINT "Hello, World!"] [PRINT "World"]]

[DEFINE !a "Hello"]
[PRINT !a]

[FUNC @hello [
    [PRINT "Hello, World!"]
]]
"#;

fn main() {
    let mut fset = FSet::new();
    let ind = fset
        .add_file_source("main.bsp".to_string(), SOURCE.to_string())
        .unwrap();
    let mut ir = IR::new(fset);
    ir.build_file(ind).unwrap();
}
