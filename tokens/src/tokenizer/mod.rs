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
            if let Some(tok) = next {
                self.tokens.push(tok)
            }
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

    pub fn peek(&self) -> Result<Token, TokenizeError> {
        if let Some(v) = self.tokens.get(self.index) {
            Ok(v.clone())
        } else {
            Err(TokenizeError::EOF)
        }
    }

    pub fn has_next(&mut self) -> bool {
        self.tokens.get(self.index).is_some()
    }

    pub fn expect(&mut self, kind: Discriminant<TokenData>) -> Result<Token, TokenizeError> {
        if let Ok(next) = self.peek() {
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
