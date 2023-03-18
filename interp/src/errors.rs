use super::*;

#[derive(Debug)]
pub enum InterpError {
    UnknownNode(IRNode),
}
