use super::*;
use fset::FSetError;

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
    pub fn pos(&self) -> Option<Pos> {
        match self {
            IRError::UnexpectedNode(node) => Some(node.pos),
            IRError::InvalidGlobalDef(node) => Some(node.pos),
            IRError::UnknownStmt { pos, .. }
            | IRError::InvalidArgumentCount { pos, .. }
            | IRError::InvalidType { pos, .. }
            | IRError::UnknownVariable { pos, .. }
            | IRError::UnknownField { pos, .. }
            | IRError::UnknownFunction { pos, .. }
            | IRError::MissingStructFields { pos, .. }
            | IRError::ReturnStatementOutsideFunction(pos)
            | IRError::StructOpOutsideDef(pos)
            | IRError::CaseOutsideMatch(pos)
            | IRError::DuplicateType(pos, _)
            | IRError::DuplicateFunction(pos, _)
            | IRError::DuplicateVariable(pos, _)
            | IRError::UnknownType { pos, .. } => Some(*pos),
            IRError::InvalidASTArgument { got, .. } => Some(got.pos),
            IRError::InvalidArgument { got, .. } => Some(got.pos),
            IRError::FSetError(e) => e.pos(),
        }
    }

    pub fn msg(&self, ir: &IR) -> String {
        match self {
            IRError::UnexpectedNode(_) => "unexpected node".to_string(),
            IRError::UnknownStmt { name, .. } => {
                format!("unknown statement '{}'", name)
            }
            IRError::InvalidGlobalDef(_) => "invalid in global scope".to_string(),
            IRError::InvalidArgumentCount { expected, got, .. } => {
                format!("invalid argument count, expected {}, got {}", expected, got)
            }
            IRError::InvalidASTArgument { expected, got } => format!(
                "invalid argument kind, expected {:?}, got {:?}",
                expected,
                got.data.typ()
            ),
            IRError::InvalidArgument { expected, got } => format!(
                "invalid argument type, expected {}, got {}",
                expected.fmt(ir),
                got.typ(ir).data.fmt(ir),
            ),
            IRError::FSetError(e) => e.msg(),
            IRError::InvalidType { expected, got, .. } => {
                format!("invalid type, expected {:?}, got {:?}", expected, got)
            }
            IRError::UnknownType { name, .. } => {
                format!("unknown type '{}'", name)
            }
            IRError::UnknownVariable { name, .. } => {
                format!("unknown variable '{}'", name)
            }
            IRError::UnknownField { name, .. } => {
                format!("unknown field '{}'", name)
            }
            IRError::UnknownFunction { name, .. } => {
                format!("unknown function '{}'", name)
            }
            IRError::MissingStructFields { missing, .. } => {
                format!("missing struct fields: {}", missing.join(", "))
            }
            IRError::ReturnStatementOutsideFunction(_) => {
                "return statement outside function".to_string()
            }
            IRError::StructOpOutsideDef(_) => "struct op outside struct definition".to_string(),
            IRError::CaseOutsideMatch(_) => "case outside match statement".to_string(),
            IRError::DuplicateType(_, orig) => {
                format!(
                    "duplicate type (originally defined at {})",
                    ir.fset.display_pos(&ir.types[*orig].pos)
                )
            }
            IRError::DuplicateFunction(_, orig) => {
                format!(
                    "duplicate function (originally defined at {})",
                    ir.fset.display_pos(&ir.funcs[*orig].definition)
                )
            }
            IRError::DuplicateVariable(_, orig) => {
                format!(
                    "duplicate variable (originally defined at {})",
                    ir.fset.display_pos(&ir.variables[*orig].definition)
                )
            }
        }
    }

    pub fn fmt(&self, ir: &IR) -> String {
        let pos = self.pos();
        if let None = pos {
            return self.msg(ir);
        }
        format!("{}: {}", ir.fset.display_pos(&pos.unwrap()), self.msg(ir))
    }
}
