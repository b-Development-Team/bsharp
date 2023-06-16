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

    pub fn build_fn(&mut self, f: usize) -> Result<Node, BStarError> {
        Ok(Node::String("UNIMPLEMENTED".to_string()))
    }

    pub fn build(&mut self) -> Result<Vec<Node>, BStarError> {
        let mut res = Vec::new();
        for i in 0..self.ir.funcs.len() {
            let v = self.build_fn(i)?;
            res.push(v);
        }
        res.push(Node::Tag("MAIN".to_string(), vec![]));
        Ok(res)
    }
}
