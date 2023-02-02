use std::mem::Discriminant;

use super::*;
mod tokenize;

pub struct Tokenizer {
    pub tokens: Vec<Token>,
    stream: Stream,
    index: usize,
}

impl Tokenizer {
    pub fn new(source: String, file: usize) -> Tokenizer {
        Tokenizer {
            stream: Stream::new(source, file),
            tokens: Vec::new(),
            index: 0,
        }
    }

    pub fn tokenize(&mut self) -> Result<(), TokenizeError> {
        while !self.stream.eof() {
            let next = self.next_token()?;
            self.tokens.push(next);
        }
        Ok(())
    }

    pub fn eat(&mut self) -> Result<Token, TokenizeError> {
        let next = self.tokens.get(self.index);
        if let Some(next) = next {
            self.index += 1;
            Ok(next.clone())
        } else {
            Err(TokenizeError::EOF)
        }
    }

    fn next(&mut self) -> Option<Token> {
        self.tokens.get(self.index).cloned()
    }

    pub fn expect(&mut self, kind: Discriminant<TokenData>) -> Result<Token, TokenizeError> {
        if let Some(next) = self.next() {
            if std::mem::discriminant(&next.data) == kind {
                self.index += 1;
                Ok(next)
            } else {
                Err(TokenizeError::UnexpectedToken {
                    got: next.clone(),
                    expected: kind,
                })
            }
        } else {
            Err(TokenizeError::EOF)
        }
    }
}
