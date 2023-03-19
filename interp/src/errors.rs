use std::io;

use super::*;

#[derive(Debug)]
pub enum InterpError {
    UnknownNode(IRNode),
    IOError(io::Error),
    InvalidEnumType { pos: Pos, expected: Type, got: Type },
}

impl From<io::Error> for InterpError {
    fn from(v: io::Error) -> Self {
        InterpError::IOError(v)
    }
}
