use super::*;
use parser::ParseError;
use tokens::TokenizeError;

pub enum FSetError {
    ParseError(ParseError),
    TokenizeError(TokenizeError),
    IOError(std::io::Error),
}

impl From<ParseError> for FSetError {
    fn from(e: ParseError) -> FSetError {
        FSetError::ParseError(e)
    }
}

impl From<TokenizeError> for FSetError {
    fn from(e: TokenizeError) -> FSetError {
        FSetError::TokenizeError(e)
    }
}

impl From<std::io::Error> for FSetError {
    fn from(e: std::io::Error) -> FSetError {
        FSetError::IOError(e)
    }
}

// Impl display for FSetError
impl FSetError {
    pub fn pos(&self) -> Option<Pos> {
        match self {
            FSetError::IOError(_) => None,
            FSetError::ParseError(e) => match e {
                ParseError::UnexpectedToken(t) => Some(t.pos),
                ParseError::EOF => None,
            },
            FSetError::TokenizeError(e) => match e {
                TokenizeError::UnexpectedToken { got: t, .. } => Some(t.pos),
                TokenizeError::EOF => None,
                TokenizeError::UnexpectedChar(_, pos) => Some(*pos),
                TokenizeError::InvalidEscapeCode(_, pos) => Some(*pos),
                TokenizeError::InvalidInt(_, pos) => Some(*pos),
                TokenizeError::InvalidFloat(_, pos) => Some(*pos),
            },
        }
    }

    pub fn msg(&self) -> String {
        match self {
            FSetError::ParseError(e) => match e {
                ParseError::UnexpectedToken(_) => "unexpected token".to_string(),
                ParseError::EOF => "parse error: un-closed bracket".to_string(),
            },
            FSetError::TokenizeError(e) => match e {
                TokenizeError::UnexpectedToken { .. } => "unexpected token".to_string(),
                TokenizeError::EOF => format!("tokenize error: unexpected EOF"),
                TokenizeError::UnexpectedChar(c, ..) => {
                    format!("invalid character '{}'", c)
                }
                TokenizeError::InvalidEscapeCode(c, ..) => {
                    format!("invalid escape code '{}'", c)
                }
                TokenizeError::InvalidInt(val, ..) => {
                    format!("invalid integer '{}'", val)
                }
                TokenizeError::InvalidFloat(val, ..) => {
                    format!("invalid float '{}'", val)
                }
            },
            FSetError::IOError(e) => e.to_string(),
        }
    }

    pub fn fmt(&self, fset: &FSet) -> String {
        match self {
            FSetError::ParseError(e) => match e {
                ParseError::UnexpectedToken(t) => {
                    format!("{}: unexpected token", fset.display_pos(&t.pos))
                }
                ParseError::EOF => format!("parse error: un-closed bracket"),
            },
            FSetError::TokenizeError(e) => match e {
                TokenizeError::UnexpectedToken { got: t, .. } => {
                    format!("{}: unexpected token", fset.display_pos(&t.pos))
                }
                TokenizeError::EOF => format!("tokenize error: unexpected EOF"),
                TokenizeError::UnexpectedChar(c, pos) => {
                    format!("{}: invalid character '{}'", fset.display_pos(&pos), c)
                }
                TokenizeError::InvalidEscapeCode(c, pos) => {
                    format!("{}: invalid escape code '{}'", fset.display_pos(&pos), c)
                }
                TokenizeError::InvalidInt(val, pos) => {
                    format!("{}: invalid integer '{}'", fset.display_pos(&pos), val)
                }
                TokenizeError::InvalidFloat(val, pos) => {
                    format!("{}: invalid float '{}'", fset.display_pos(&pos), val)
                }
            },
            FSetError::IOError(e) => e.to_string(),
        }
    }
}
