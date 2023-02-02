use super::*;

impl super::Tokenizer {
    pub fn next_token(&mut self) -> Result<Option<Token>, TokenizeError> {
        let start = self.stream.pos();
        match self.stream.eat() {
            None => Err(TokenizeError::EOF),
            Some(c) => match c {
                '\n' | '\t' | ' ' => Ok(None),
                _ => Err(TokenizeError::UnexpectedChar(c, self.stream.pos())),
            },
        }
    }

    fn parse_string(&mut self, mut pos: Pos) -> Token {
        let val = String::new();
        panic!("unimplemented");
    }
}
