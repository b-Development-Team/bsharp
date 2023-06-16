use ir::*;
mod node;
mod stmts;
use node::*;

pub struct BStar {
    ir: IR,
}

impl BStar {
    pub fn new(ir: IR) -> BStar {
        BStar { ir }
    }
}
