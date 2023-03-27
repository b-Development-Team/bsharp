use super::*;
use std::mem::Discriminant;

pub enum TokenizeError {
    UnexpectedToken {
        got: Token,
        expected: Discriminant<TokenData>,
    },
    UnexpectedChar(char, Pos),
    InvalidEscapeCode(char, Pos),
    InvalidInt(String, Pos),
    InvalidFloat(String, Pos),
    EOF,
}
