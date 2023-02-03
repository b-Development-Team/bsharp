use parser::ParseError;
use tokens::TokenizeError;

#[derive(Debug)]
pub enum FSetError {
    ParseError(ParseError),
    TokenizeError(TokenizeError),
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
