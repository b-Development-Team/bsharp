use super::*;
use std::mem::Discriminant;

#[derive(Debug)]
pub enum TokenizeError {
    UnexpectedToken {
        got: Token,
        expected: Discriminant<TokenData>,
    },
    UnexpectedChar(char, Pos),
    EOF,
}
