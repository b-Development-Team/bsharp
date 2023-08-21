use super::*;

pub struct Stream {
    source: Vec<char>,
    index: usize,

    file: usize,
    line: usize,
    col: usize,
}

impl Stream {
    pub fn new(source: String, file: usize) -> Stream {
        Stream {
            source: source.chars().collect(),
            index: 0,
            line: 0,
            col: 0,
            file,
        }
    }

    pub fn peek(&self) -> Option<char> {
        if self.index < self.source.len() {
            Some(self.source[self.index])
        } else {
            None
        }
    }

    pub fn eat(&mut self) -> Option<char> {
        let next = self.peek();
        self.index += 1;
        if let Some(v) = next {
            if v == '\n' {
                self.line += 1;
                self.col = 0;
            } else {
                self.col += 1;
            }
        } else {
            self.col += 1;
        }
        next
    }

    pub fn eof(&self) -> bool {
        self.peek().is_none()
    }

    pub fn pos(&self) -> Pos {
        Pos {
            file: self.file,
            start_line: self.line,
            start_col: self.col,
            end_line: self.line,
            end_col: self.col,
        }
    }

    pub fn last_char(&self) -> Pos {
        Pos {
            file: self.file,
            start_line: self.line,
            start_col: self.col - 1,
            end_line: self.line,
            end_col: self.col,
        }
    }
}

impl Pos {
    pub fn extend(&self, new: Pos) -> Pos {
        Pos {
            file: self.file,
            start_line: self.start_line,
            start_col: self.start_col,
            end_line: new.start_line,
            end_col: new.start_col,
        }
    }
}
