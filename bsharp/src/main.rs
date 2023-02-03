use parser;
use tokens;

const SOURCE: &'static str = r#"
# Block comment
# Small comment #
[PRINT "Hello, World!"]
[+ 1 [- 2 1] [% 1 2]]

# Block
[[PRINT "Hello, World!"] [PRINT "World"]]
"#;

fn main() {
    let mut tok = tokens::Tokenizer::new(SOURCE.to_string(), 0);
    tok.tokenize().unwrap();
    let mut p = parser::Parser::new(tok);
    p.parse().unwrap();
    println!("{:#?}", p.ast);
}
