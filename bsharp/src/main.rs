use clap::{error::ErrorKind, CommandFactory, Parser, ValueEnum};
use fset::FSet;
use ir::IR;

mod run;
use run::run;
mod bstar_run;
use bstar_run::bstar;

#[derive(ValueEnum, Clone, Copy)]
enum Action {
    Run,
    Bstar,
}

#[derive(Parser)]
#[command(name = "bsharp", about = "The B# programming language!")]
struct Args {
    #[arg(value_enum, help = "the action to perform on the code")]
    action: Action,

    #[arg(help = "the directory to import")]
    dir: std::path::PathBuf,
}

fn main() {
    let args = Args::parse();

    // Fset
    let mut fset = FSet::new();
    if let Err(err) = fset.import(&args.dir) {
        Args::command().error(ErrorKind::Io, err.fmt(&fset)).exit();
    };

    // IR
    let mut ir = IR::new(fset);
    if let Err(err) = ir.build() {
        Args::command()
            .error(ErrorKind::InvalidValue, err.fmt(&ir))
            .exit();
    };
    if ir.errors.len() > 0 {
        for err in ir.errors.iter() {
            println!("{}", err.fmt(&ir));
        }
        println!("");
        Args::command()
            .error(
                ErrorKind::ValueValidation,
                format!("built with {} errors", ir.errors.len()),
            )
            .exit();
    }

    // Run
    match args.action {
        Action::Run => run(ir),
        Action::Bstar => bstar(ir),
    }
}
