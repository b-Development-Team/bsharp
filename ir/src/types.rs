#[derive(Debug)]
pub struct Type {
    pub data: TypeData,
    pub name: Option<String>,
}

impl Type {
    pub fn void() -> Self {
        Self {
            data: TypeData::VOID,
            name: None,
        }
    }
}

#[derive(Debug)]
pub enum TypeData {
    INT,
    FLOAT,
    CHAR,
    BOOL,
    BOX,
    ARRAY(Box<Type>),
    STRUCT {
        params: Vec<Generic>,
        fields: Vec<Field>,
    },
    TUPLE {
        params: Vec<Generic>,
        body: Vec<Type>,
    },
    ENUM {
        params: Vec<Generic>,
        body: Vec<Type>,
    },
    INTERFACE {
        params: Vec<Generic>,
        body: Vec<Type>,
    },

    // Special types
    INVALID,
    TYPE,
    VOID,
}

#[derive(Debug)]
pub struct Generic {
    pub name: String,
    pub typ: Type, // INTERFACE
}

#[derive(Debug)]
pub struct Field {
    pub name: String,
    pub typ: Type,
}
