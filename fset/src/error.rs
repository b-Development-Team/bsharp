use parser::ParseError;
use tokens::TokenizeError;

#[derive(Debug)]
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
