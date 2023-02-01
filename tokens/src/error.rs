use std::mem::Discriminant;
use super::*;

#[derive(Debug)]
pub enum TokenizeError {
  UnexpectedToken{
    got: Token,
    expected: Discriminant<TokenData>,
  },
  EOF,
}