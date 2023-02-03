use super::*;

impl super::Tokenizer {
    pub fn next_token(&mut self) -> Result<Option<Token>, TokenizeError> {
        match self.stream.eat() {
            None => Err(TokenizeError::EOF),
            Some(c) => match c {
                '"' => Ok(Some(self.parse_string()?)),
                '\n' | '\t' | ' ' => Ok(None),
                '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' => {
                    Ok(Some(self.parse_num(c)?))
                }
                '\'' => Ok(Some(self.parse_char()?)),
                _ => Err(TokenizeError::UnexpectedChar(c, self.stream.pos())),
            },
        }
    }

    fn parse_char(&mut self) -> Result<Token, TokenizeError> {
        let pos = self.stream.last_char();
        let val = match self.stream.eat() {
            None => return Err(TokenizeError::EOF),
            Some('\\') => match self.stream.eat() {
                None => return Err(TokenizeError::EOF),
                Some('n') => '\n',
                Some('t') => '\t',
                Some('\'') => '\'',
                Some('\\') => '\\',
                Some(v) => {
                    return Err(TokenizeError::InvalidEscapeCode(v, self.stream.last_char()))
                }
            },
            Some(v) => v,
        };
        match self.stream.eat() {
            Some('\'') => {}
            Some(v) => return Err(TokenizeError::UnexpectedChar(v, self.stream.last_char())),
            None => return Err(TokenizeError::EOF),
        }
        return Ok(Token {
            data: TokenData::CHAR(val),
            pos: pos.extend(self.stream.pos()),
        });
    }

    fn parse_string(&mut self) -> Result<Token, TokenizeError> {
        let pos = self.stream.last_char();
        let mut val = String::new();
        loop {
            let c = self.stream.eat();
            if let Some(v) = c {
                match v {
                    '\\' => val.push(match self.stream.eat() {
                        None => return Err(TokenizeError::EOF),
                        Some('n') => '\n',
                        Some('t') => '\t',
                        Some('"') => '"',
                        Some('\\') => '\\',
                        Some(v) => {
                            return Err(TokenizeError::InvalidEscapeCode(
                                v,
                                self.stream.last_char(),
                            ))
                        }
                    }),
                    '"' => {
                        return Ok(Token {
                            data: TokenData::STRING(val),
                            pos: pos.extend(self.stream.pos()),
                        });
                    }
                    _ => val.push(v),
                }
            } else {
                return Err(TokenizeError::EOF);
            }
        }
    }

    fn parse_num(&mut self, start: char) -> Result<Token, TokenizeError> {
        let mut val = start.to_string();
        let mut pos = self.stream.last_char();
        let mut int = true;
        loop {
            let next = self.stream.peek();
            if let Some(v) = next {
                match v {
                    '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' => {
                        val.push(self.stream.eat().unwrap())
                    }
                    '.' => {
                        if int {
                            int = false;
                        }
                        val.push(self.stream.eat().unwrap());
                    }
                    _ => {
                        pos = pos.extend(self.stream.pos());
                        let dat = if int {
                            let res = val.parse::<i64>();
                            if let Ok(v) = res {
                                TokenData::INTEGER(v)
                            } else {
                                return Err(TokenizeError::InvalidInt(val, pos));
                            }
                        } else {
                            let res = val.parse::<f64>();
                            if let Ok(v) = res {
                                TokenData::FLOAT(v)
                            } else {
                                return Err(TokenizeError::InvalidFloat(val, pos));
                            }
                        };
                        return Ok(Token { data: dat, pos });
                    }
                }
            } else {
                return Err(TokenizeError::EOF);
            }
        }
    }
}
