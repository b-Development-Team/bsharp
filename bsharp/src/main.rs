use tokens;

const SOURCE: &'static str = r#"
[PRINT "Hello, World!"]
"#;

fn main() {
    let mut tok = tokens::Tokenizer::new(SOURCE.to_string());
    tok.tokenize().unwrap();
    println!("{:?}", tok.tokens);
}
