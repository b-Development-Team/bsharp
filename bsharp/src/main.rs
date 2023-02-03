use fset::FSet;

const SOURCE: &'static str = r#"
# Block comment
# Small comment #
[PRINT "Hello, World!"]
[+ 1 [- 2 1] [% 1 2]]

# Block
[[PRINT "Hello, World!"] [PRINT "World"]]
"#;

fn main() {
    let mut fset = FSet::new();
    let ind = fset
        .add_file_source("main.bsp".to_string(), SOURCE.to_string())
        .unwrap();
    let f = fset.get_file(ind).unwrap();
    println!("{:#?}",f.ast.ast);
}
