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
    Stmt {
        name: String,
        name_pos: Pos,
        args: Vec<ASTNode>,
    },
    Block(Vec<ASTNode>),
}
