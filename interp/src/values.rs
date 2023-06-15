use std::rc::Rc;
use std::{cell::RefCell, collections::HashMap};

use super::*;

#[derive(Clone, Debug, PartialEq)]
pub enum Value {
    Int(i64),
    Float(f64),
    Char(u8),
    Bool(bool),
    Box(Type, Box<Value>),
    Array(Rc<RefCell<Vec<Value>>>),
    Struct(Rc<RefCell<HashMap<String, Value>>>),
    Tuple(Rc<RefCell<Vec<Value>>>),
    Enum(Type, Box<Value>),
    Void,
}

#[derive(Default)]
pub struct StackFrame {
    pub vars: HashMap<usize, Value>,
    pub ret: Option<Value>,
    pub pos: Pos,
    pub func: usize,
}
