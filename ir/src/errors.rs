use super::*;

pub enum IRError {
    UnexpectedNode(ASTNode),
    UnknownStmt {
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
    UnknownType {
        pos: Pos,
        name: String,
    },
    UnknownVariable {
        pos: Pos,
        name: String,
    },
    UnknownField {
        pos: Pos,
        name: String,
    },
    UnknownFunction {
        pos: Pos,
        name: String,
    },
    MissingStructFields {
        pos: Pos,
        missing: Vec<String>,
    },
    ReturnStatementOutsideFunction(Pos),
    StructOpOutsideDef(Pos),
    CaseOutsideMatch(Pos),
    DuplicateType(Pos, usize),     // Usize has original type index
    DuplicateFunction(Pos, usize), // Usize has original function index
    DuplicateVariable(Pos, usize), // Usize has original variable index
}

impl From<FSetError> for IRError {
    fn from(e: FSetError) -> Self {
        IRError::FSetError(e)
    }
}

impl IRError {
    pub fn fmt(&self, ir: &IR) -> String {
        match self {
            IRError::UnexpectedNode(node) => {
                format!("{}: unexpected node", ir.fset.display_pos(&node.pos))
            }
            IRError::UnknownStmt { pos, name } => {
                format!("{}: unknown statement '{}'", ir.fset.display_pos(pos), name)
            }
            IRError::InvalidGlobalDef(node) => {
                format!(
                    "{}: invalid in global scope",
                    ir.fset.display_pos(&node.pos)
                )
            }
            IRError::InvalidArgumentCount { pos, expected, got } => format!(
                "{}: invalid argument count, expected {}, got {}",
                ir.fset.display_pos(pos),
                expected,
                got
            ),
            IRError::InvalidASTArgument { expected, got } => format!(
                "{}: invalid argument kind, expected {:?}, got {:?}",
                ir.fset.display_pos(&got.pos),
                expected,
                got.data.typ()
            ),
            IRError::InvalidArgument { expected, got } => format!(
                "{}: invalid argument type, expected {:?}, got {:?}",
                ir.fset.display_pos(&got.pos),
                expected,
                got.typ(ir),
            ),
            IRError::FSetError(e) => format!("{}", e),
            IRError::InvalidType { pos, expected, got } => format!(
                "{}: invalid type, expected {:?}, got {:?}",
                ir.fset.display_pos(pos),
                expected,
                got
            ),
            IRError::UnknownType { pos, name } => {
                format!("{}: unknown type '{}'", ir.fset.display_pos(pos), name)
            }
            IRError::UnknownVariable { pos, name } => {
                format!("{}: unknown variable '{}'", ir.fset.display_pos(pos), name)
            }
            IRError::UnknownField { pos, name } => {
                format!("{}: unknown field '{}'", ir.fset.display_pos(pos), name)
            }
            IRError::UnknownFunction { pos, name } => {
                format!("{}: unknown function '{}'", ir.fset.display_pos(pos), name)
            }
            IRError::MissingStructFields { pos, missing } => {
                format!(
                    "{}: missing struct fields: {}",
                    ir.fset.display_pos(pos),
                    missing.join(", ")
                )
            }
            IRError::ReturnStatementOutsideFunction(pos) => {
                format!(
                    "{}: return statement outside function",
                    ir.fset.display_pos(pos)
                )
            }
            IRError::StructOpOutsideDef(pos) => {
                format!(
                    "{}: struct op outside struct definition",
                    ir.fset.display_pos(pos)
                )
            }
            IRError::CaseOutsideMatch(pos) => {
                format!("{}: case outside match statement", ir.fset.display_pos(pos))
            }
            IRError::DuplicateType(pos, orig) => {
                format!(
                    "{}: duplicate type (originally defined at {})",
                    ir.fset.display_pos(pos),
                    ir.fset.display_pos(&ir.types[*orig].pos)
                )
            }
            IRError::DuplicateFunction(pos, orig) => {
                format!(
                    "{}: duplicate function (originally defined at {})",
                    ir.fset.display_pos(pos),
                    ir.fset.display_pos(&ir.funcs[*orig].definition)
                )
            }
            IRError::DuplicateVariable(pos, orig) => {
                format!(
                    "{}: duplicate variable (originally defined at {})",
                    ir.fset.display_pos(pos),
                    ir.fset.display_pos(&ir.variables[*orig].definition)
                )
            }
        }
    }
}
