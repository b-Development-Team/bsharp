use std::io;

use super::*;

#[derive(Debug)]
pub enum InterpError {
    UnknownNode(IRNode),
    IOError(io::Error),
    InvalidEnumType { pos: Pos, expected: Type, got: Type },
    InvalidBoxType { pos: Pos, expected: Type, got: Type },
    ArrayIndexOutOfBounds { pos: Pos, index: usize, len: usize },
    DivideByZero(Pos),
}

impl From<io::Error> for InterpError {
    fn from(v: io::Error) -> Self {
        InterpError::IOError(v)
    }
}

impl InterpError {
    pub fn fmt(&self, int: &Interp) -> String {
        match self {
            InterpError::UnknownNode(node) => format!("unknown node: {:?}", node),
            InterpError::IOError(e) => format!("IO error: {}", e),
            InterpError::InvalidEnumType { pos, expected, got } => format!(
                "{}: invalid enum type, expected {}, got {}",
                int.ir.fset.display_pos(pos),
                expected.data.fmt(&int.ir),
                got.data.fmt(&int.ir),
            ),
            InterpError::InvalidBoxType { pos, expected, got } => format!(
                "{}: invalid box type, expected {}, got {}",
                int.ir.fset.display_pos(pos),
                expected.data.fmt(&int.ir),
                got.data.fmt(&int.ir),
            ),
            InterpError::ArrayIndexOutOfBounds { pos, index, len } => format!(
                "{}: array index out of bounds, index {} >= len {}",
                int.ir.fset.display_pos(pos),
                index,
                len,
            ),
            InterpError::DivideByZero(pos) => {
                format!("{}: divide by zero", int.ir.fset.display_pos(pos),)
            }
        }
    }
}
