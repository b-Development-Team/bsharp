mod types;
use fset::*;
pub use types::*;
mod stack;
use stack::*;
mod ir;
pub use ir::*;

pub struct IR {
    pub scopes: Vec<Scope>,
    pub variables: Vec<Variable>,
    pub types: Vec<TypeDef>,

    fset: FSet,
    stack: Vec<usize>,
}

impl IR {
    pub fn new() -> Self {
        Self {
            fset: FSet::new(),
            variables: Vec::new(),
            types: Vec::new(),
            scopes: Vec::new(),
            stack: Vec::new(),
        }
    }
}
