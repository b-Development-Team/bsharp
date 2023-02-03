use tokens::Tokenizer;
mod error;
use error::FSetError;

pub struct File {
    pub name: String,
    pub ast: parser::Parser,
}

pub struct FSet {
    files: Vec<File>,
}

impl FSet {
    pub fn new() -> FSet {
        FSet { files: Vec::new() }
    }

    pub fn add_file_source(&mut self, name: String, source: String) -> Result<usize, FSetError> {
        let ind = self.files.len();
        let mut tok = Tokenizer::new(source, ind);
        tok.tokenize()?;

        let mut p = parser::Parser::new(tok);
        p.parse()?;

        self.files.push(File { name, ast: p });

        Ok(ind)
    }

    pub fn get_file(&self, ind: usize) -> Option<&File> {
        self.files.get(ind)
    }
}
