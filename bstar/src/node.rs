use super::*;
use std::fmt;

#[derive(Clone)]
pub enum Node {
    String(String),
    Ident(String),
    Int(i64),
    Float(f64),
    Tag(String, Vec<Node>),
    ArrayLiteral(Vec<Node>),
}

impl fmt::Display for Node {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            Self::String(v) => write!(f, "{:?}", v),
            Self::Ident(v) => write!(f, "{}", v),
            Self::Int(v) => write!(f, "{}", v),
            Self::Float(v) => write!(f, "{}", v),
            Self::Tag(name, vals) => {
                write!(f, "[{}", name)?;
                for val in vals {
                    write!(f, " {}", val)?;
                }
                write!(f, "]")
            }
            Self::ArrayLiteral(vals) => {
                write!(f, "{{")?;
                for (i, val) in vals.iter().enumerate() {
                    if i != 0 {
                        write!(f, ", ")?;
                    }
                    write!(f, "{}", val)?;
                }
                write!(f, "}}")
            }
        }
    }
}

#[derive(Debug)]
pub enum BStarError {
    UnknownNode(IRNode),
}
