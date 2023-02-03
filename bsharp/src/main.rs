use tokens;

const SOURCE: &'static str = r#"
"Hello\n \"World\""
1
1.3
'1'
'\''
IDENT
# Block comment
# Small comment #
[PRINT "Hello, World!"]

# Block
[[PRINT "Hello, World!"] [PRINT "World"]]
"#;

fn main() {
    let mut tok = tokens::Tokenizer::new(SOURCE.to_string(), 0);
    tok.tokenize().unwrap();
    println!("{:?}", tok.tokens);
}
