use super::*;
use std::fmt;

pub enum Node {
    String(String),
    Int(i64),
    Float(f64),
    Tag(String, Vec<Node>),
}

impl fmt::Display for Node {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Self::String(v) => write!(f, "{:?}", v),
            Self::Int(v) => write!(f, "{}", v),
            Self::Float(v) => write!(f, "{}", v),
            Self::Tag(name, vals) => {
                write!(f, "[{}", name)?;
                for val in vals {
                    write!(f, " {}", val)?;
                }
                write!(f, "]")
            }
        }
    }
}

#[derive(Debug)]
pub enum BStarError {
    UnknownNode(IRNode),
}
