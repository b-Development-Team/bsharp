use super::*;

#[derive(Debug)]
pub enum ParseError {
    UnexpectedToken(Token),
    EOF,
}

impl From<TokenizeError> for ParseError {
    fn from(e: TokenizeError) -> Self {
        match e {
            TokenizeError::UnexpectedToken {
                got: t,
                expected: _,
            } => ParseError::UnexpectedToken(t),
            TokenizeError::EOF => ParseError::EOF,
            _ => unreachable!(),
        }
    }
}
