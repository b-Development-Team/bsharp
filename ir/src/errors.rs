use super::*;

#[derive(Debug)]
pub enum IRError {
    UnexpectedNode(ASTNode),
    InvalidGlobalDef(IRNode),
}
