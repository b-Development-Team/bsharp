use super::*;
use interp::Interp;

pub fn run(ir: IR) {
    let mut interp = Interp::new(ir, Box::new(std::io::stdout()));
    let mainind = interp.ir.funcs.iter().position(|f| f.name == "@MAIN");
    if mainind.is_none() {
        panic!("main func not found");
    }
    if let Err(err) = interp.run_fn(
        mainind.unwrap(),
        Vec::new(),
        interp.ir.funcs[mainind.unwrap()].definition,
    ) {
        // Print stack trace
        for frame in interp.stack.iter().rev() {
            println!(
                "{} at {}",
                interp.ir.funcs[frame.func].name,
                interp.ir.fset.display_pos(&frame.pos)
            );
        }

        Args::command()
            .error(ErrorKind::InvalidValue, err.fmt(&interp))
            .exit();
    };
}
