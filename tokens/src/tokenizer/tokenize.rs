use super::*;

impl super::Tokenizer {
    pub fn next_token(&mut self) -> Result<Token, TokenizeError> {
        match self.stream.eat() {
            None => Err(TokenizeError::EOF),
            Some(c) => match c {
                _ => Err(TokenizeError::UnexpectedChar(c, self.stream.pos())),
            },
        }
    }
}
