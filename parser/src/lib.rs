use tokens::{Pos, Token, TokenData, TokenizeError, Tokenizer};
mod ast;
mod error;
mod parse;

pub use ast::*;
pub use error::ParseError;

pub struct Parser {
    tok: Tokenizer,
    pub ast: Vec<Node>, // TODO: Make private
}

impl Parser {
    // new accepts a tokenizer that has already been tokenized
    pub fn new(tok: Tokenizer) -> Self {
        Self {
            tok,
            ast: Vec::new(),
        }
    }
}
