#[derive(Debug, Clone)]
pub struct Type {
    pub data: TypeData,
    pub name: Option<String>,
}

impl PartialEq for Type {
    fn eq(&self, other: &Self) -> bool {
        self.data == other.data
    }
}

impl From<TypeData> for Type {
    fn from(data: TypeData) -> Self {
        Self { data, name: None }
    }
}

impl Type {
    pub fn void() -> Self {
        Self {
            data: TypeData::VOID,
            name: None,
        }
    }
}

#[derive(Debug, PartialEq, Clone)]
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

#[derive(Debug, PartialEq, Clone)]
pub struct Generic {
    pub name: String,
    pub typ: Type, // INTERFACE
}

#[derive(Debug, PartialEq, Clone)]
pub struct Field {
    pub name: String,
    pub typ: Type,
}
