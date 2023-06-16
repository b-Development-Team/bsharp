use super::*;
use bstar::BStar;

pub fn bstar(ir: IR) {
    let mut bstar = BStar::new(ir);
    let res = bstar.build().unwrap(); // TODO: Error handling
    for i in res {
        println!("{}", i);
    }
}
