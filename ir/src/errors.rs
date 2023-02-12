use super::*;

#[derive(Debug)]
pub enum IRError {
    UnexpectedNode(ASTNode),
    InvalidGlobalDef(IRNode),
    InvalidArgumentCount {
        pos: Pos,
        expected: usize,
        got: usize,
    },
    InvalidASTArgument {
        expected: ASTNodeDataType,
        got: ASTNode,
    },
    FSetError(FSetError),
}

impl From<FSetError> for IRError {
    fn from(e: FSetError) -> Self {
        IRError::FSetError(e)
    }
}
