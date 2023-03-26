use std::io;

use super::*;

#[derive(Debug)]
pub enum InterpError {
    UnknownNode(IRNode),
    IOError(io::Error),
    InvalidEnumType { pos: Pos, expected: Type, got: Type },
    InvalidBoxType { pos: Pos, expected: Type, got: Type },
    ArrayIndexOutOfBounds { pos: Pos, index: usize, len: usize },
}

impl From<io::Error> for InterpError {
    fn from(v: io::Error) -> Self {
        InterpError::IOError(v)
    }
}
