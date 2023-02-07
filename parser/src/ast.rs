use super::*;

#[derive(Debug)]
pub struct Node {
    pub data: NodeData,
    pub pos: Pos,
}

#[derive(Debug)]
pub enum NodeData {
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
        args: Vec<Node>,
    },
    Block(Vec<Node>),
}
