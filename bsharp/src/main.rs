use tokens;

const SOURCE: &'static str = r#"
# Block comment
# Small comment #
[PRINT "Hello, World!"]

# Block
[[PRINT "Hello, World!"] [PRINT "World"]]
[+ 1 [- 2 1] [% 1 2]]
"#;

fn main() {
    let mut tok = tokens::Tokenizer::new(SOURCE.to_string(), 0);
    tok.tokenize().unwrap();
    println!("{:?}", tok.tokens);
}
