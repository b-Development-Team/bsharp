#[derive(Clone, Debug)]
pub struct Token {
  pub data: TokenData,
  pub pos: Pos,
}

#[derive(Clone, Copy, Debug)]
pub struct Pos {
  pub file: usize,
  pub line: usize,
  pub col: usize,
}

#[derive(Clone, Debug)]
pub enum TokenData {
  IDENT(String),
}