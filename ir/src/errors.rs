use super::*;

#[derive(Debug)]
pub enum IRError {
    UnexpectedNode(ASTNode),
    UnknownFunction {
        pos: Pos,
        name: String,
    },
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
    InvalidArgument {
        expected: TypeData,
        got: IRNode,
    },
    FSetError(FSetError),
    InvalidType {
        pos: Pos,
        expected: Type,
        got: Type,
    },
}

impl From<FSetError> for IRError {
    fn from(e: FSetError) -> Self {
        IRError::FSetError(e)
    }
}
