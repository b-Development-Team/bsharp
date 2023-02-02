use tokens;

const SOURCE: &'static str = r#"
"Hello\n \"World\""
1
1.3
'1'
'\''
[PRINT "Hello, World!"]
"#;

fn main() {
    let mut tok = tokens::Tokenizer::new(SOURCE.to_string(), 0);
    tok.tokenize().unwrap();
    println!("{:?}", tok.tokens);
}
