use super::*;

impl BStar {
    pub fn build_node(&mut self, node: &IRNode) -> Result<Node, BStarError> {
        match node {
            _ => Err(BStarError::UnknownNode(node.clone())),
        }
    }
}
