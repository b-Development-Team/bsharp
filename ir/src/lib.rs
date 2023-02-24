mod fns;
mod types;
use fset::*;
pub use types::*;
mod stack;
use stack::*;
mod ir;
pub use ir::*;
mod errors;
mod stmts;
use errors::*;
mod defpass;
mod typecheck;
use typecheck::*;

pub struct IR {
    pub scopes: Vec<Scope>,
    pub variables: Vec<Variable>,
    pub types: Vec<TypeDef>,
    pub funcs: Vec<Function>,

    pub errors: Vec<IRError>,

    fset: FSet,
    stack: Vec<usize>,
}

impl IR {
    pub fn new(fset: FSet) -> Self {
        Self {
            fset,
            variables: Vec::new(),
            types: Vec::new(),
            scopes: vec![Scope::new(ScopeKind::Global, Pos::default())],
            stack: vec![0],
            funcs: Vec::new(),
            errors: Vec::new(),
        }
    }

    fn save_error(&mut self, err: IRError) {
        self.errors.push(err);
    }

    pub fn build(&mut self) -> Result<(), IRError> {
        self.defpass()?;
        for i in 0..self.types.len() {
            self.build_typ(i);
        }
        for i in 0..self.funcs.len() {
            self.build_fn(i);
        }
        Ok(())
    }
}
