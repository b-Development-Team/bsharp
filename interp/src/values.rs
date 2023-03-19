use std::collections::HashMap;

use super::*;

#[derive(Clone, Debug)]
pub enum Value {
    Int(i64),
    Float(f64),
    Char(u8),
    Bool(bool),
    Box(Box<Value>),
    Array(Vec<Value>),
    Struct(HashMap<String, Value>),
    Tuple(Vec<Value>),
    Enum(Type, Box<Value>),
    Void,
}

#[derive(Default)]
pub struct StackFrame {
    pub vars: HashMap<usize, Value>,
    pub ret: Option<Value>,
}
