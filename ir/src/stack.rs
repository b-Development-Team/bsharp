use std::collections::HashMap;

use super::*;

pub struct Scope {
    pub kind: StackScopeKind,
    pub vars: HashMap<String, usize>,
    pub types: HashMap<String, usize>,
    pub pos: Pos,
}

pub struct Variable {
    pub name: String,
    pub typ: Type,
    pub scope: usize,
    pub definition: Pos,
}

pub struct TypeDef {
    pub scope: usize,
    pub name: String,
    pub typ: Type,
    pub definition: Pos,
}

pub enum StackScopeKind {
    Global,
    Type,
    Function,
    Block,
}
