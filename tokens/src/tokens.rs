#[derive(Clone, Debug)]
pub struct Token {
    pub data: TokenData,
    pub pos: Pos,
}

#[derive(Clone, Copy, Debug)]
pub struct Pos {
    pub file: usize,

    pub start_line: usize,
    pub start_col: usize,

    pub end_line: usize,
    pub end_col: usize,
}

#[derive(Clone, Debug, PartialEq)]
pub enum TokenData {
    IDENT(String),
    STRING(String),
    TYPE(String),
    INTEGER(i64),
    FLOAT(f64),
    CHAR(char),
    OPENBRACK,
    CLOSEBRACK,
    COMMENT(String),
}
