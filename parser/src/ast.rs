use super::*;

#[derive(Debug, Clone)]
pub struct ASTNode {
    pub data: ASTNodeData,
    pub pos: Pos,
}

#[derive(Debug, Clone)]
pub enum ASTNodeData {
    String(String),
    Integer(i64),
    Float(f64),
    Char(char),
    Bool(bool),
    Comment(String),
    Type(String),
    Function(String),
    Variable(String),
    Field(String),
    Stmt {
        name: String,
        name_pos: Pos,
        args: Vec<ASTNode>,
    },
    Block(Vec<ASTNode>),
}

impl ASTNodeData {
    pub fn typ(&self) -> ASTNodeDataType {
        match self {
            Self::String(_) => ASTNodeDataType::String,
            Self::Integer(_) => ASTNodeDataType::Integer,
            Self::Float(_) => ASTNodeDataType::Float,
            Self::Char(_) => ASTNodeDataType::Char,
            Self::Bool(_) => ASTNodeDataType::Bool,
            Self::Comment(_) => ASTNodeDataType::Comment,
            Self::Type(_) => ASTNodeDataType::Type,
            Self::Function(_) => ASTNodeDataType::Function,
            Self::Variable(_) => ASTNodeDataType::Variable,
            Self::Stmt { .. } => ASTNodeDataType::Stmt,
            Self::Block(_) => ASTNodeDataType::Block,
            Self::Field(_) => ASTNodeDataType::Field,
        }
    }
}

#[derive(Clone, Copy, PartialEq, Debug)]
pub enum ASTNodeDataType {
    String,
    Integer,
    Float,
    Char,
    Bool,
    Comment,
    Type,
    Function,
    Variable,
    Stmt,
    Block,
    Field,

    Any,
}
