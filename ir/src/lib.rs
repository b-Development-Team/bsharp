mod types;
use std::collections::HashMap;

use fset::*;
pub use types::*;
mod stack;
use stack::*;
mod ir;
pub use ir::*;
mod errors;
mod stmts;
use errors::*;

pub struct IR {
    pub scopes: Vec<Scope>,
    pub variables: Vec<Variable>,
    pub types: Vec<TypeDef>,
    pub funcs: Vec<Function>,

    fset: FSet,
    stack: Vec<usize>,
}

impl IR {
    pub fn new(fset: FSet) -> Self {
        Self {
            fset,
            variables: Vec::new(),
            types: Vec::new(),
            scopes: vec![Scope {
                kind: ScopeKind::Global,
                vars: HashMap::new(),
                types: HashMap::new(),
                pos: Pos {
                    file: usize::MAX,
                    start_line: 0,
                    start_col: 0,
                    end_line: 0,
                    end_col: 0,
                },
            }],
            stack: vec![0],
            funcs: Vec::new(),
        }
    }

    pub fn build_file(&mut self, pos: usize) -> Result<(), IRError> {
        let ast = self.fset.get_file(pos).unwrap().ast.ast.clone(); // TODO: Don't have to clone
        for node in ast.iter() {
            let n = self.build_node(node)?;
            match n.data {
                IRNodeData::Void => {}
                _ => return Err(IRError::InvalidGlobalDef(n)),
            }
        }
        Ok(())
    }
}
