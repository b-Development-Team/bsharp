use super::*;

#[derive(Debug)]
pub struct IRNode {
    pub data: IRNodeData,
    pub pos: Pos,
}

impl IRNode {
    pub fn void() -> Self {
        Self {
            data: IRNodeData::Void,
            pos: Pos::default(),
        }
    }
}

impl IRNode {
    pub fn new(data: IRNodeData, pos: Pos) -> Self {
        Self { data, pos }
    }
}

#[derive(Debug)]
pub enum IRNodeData {
    Block {
        scope: usize,
        body: Vec<IRNode>,
    },
    Print(Vec<IRNode>),
    Math(Box<IRNode>, MathOperator, Box<IRNode>),
    Comparison(Box<IRNode>, ComparisonOperator, Box<IRNode>),
    Boolean(Box<IRNode>, BooleanOperator, Box<IRNode>),
    Define {
        var: usize,
        val: Box<IRNode>,
        edit: bool, // Whether making variable or just changing value
    },
    Import(String),
    While {
        cond: Box<IRNode>,
        body: Box<IRNode>,
    },
    TypeCase {
        var: usize,
        typ: Type,
        body: Box<IRNode>,
    },
    Case {
        val: Box<IRNode>,
    },
    TypeMatch {
        val: Box<IRNode>,
        body: Vec<IRNode>, // TypeCase
    },
    Match {
        val: Box<IRNode>,
        body: Vec<IRNode>, // Case
    },
    If {
        cond: Box<IRNode>,
        body: Box<IRNode>, // Doesn't have to be block, can be used as ternary op
        els: Box<IRNode>,
        ret_typ: Type,
    },
    Void,

    // Composite type ops
    Len(Box<IRNode>),
    Append {
        arr: Box<IRNode>,
        val: Box<IRNode>,
    },
    GetEnum {
        enm: Box<IRNode>,
        typ: Type,
    },
    GetStruct {
        strct: Box<IRNode>,
        field: String,
    },
    StructOp {
        field: String,
        val: Box<IRNode>,
    }, // [:]
    SetStruct {
        strct: Box<IRNode>,
        vals: Box<IRNode>,
    },
    Peek {
        bx: Box<IRNode>,
        typ: Type,
    },
    Unbox {
        bx: Box<IRNode>,
        typ: Type,
    },

    // Allocating
    NewArray(Type, Option<usize>), // optional: capacity
    NewEnum(Box<IRNode>),
    NewBox(Box<IRNode>),          // [BOX]
    NewStruct(Type, Vec<IRNode>), // [:] statements

    // Functions
    Param {
        name: String,
        typ: Type,
    },
    Returns(Type),

    // Types
    Type(Type), // [INT], [STRUCT], etc.
    Generic {
        name: String,
        typ: Type,
    }, // Type parameters
    TypeInstantiate {
        typ: Type,
        params: Vec<Type>,
    }, // [G <Type> <Params>]
    Cast(Box<IRNode>, Type),
}

#[derive(Debug)]
pub enum MathOperator {
    ADD,
    SUBTRACT,
    MULTIPLY,
    DIVIDE,
    MODULO,
}

#[derive(Debug)]
pub enum ComparisonOperator {
    GREATER,
    LESS,
    GREATEREQUAL,
    LESSEQUAL,
}

#[derive(Debug)]
pub enum BooleanOperator {
    AND,
    OR,
}
