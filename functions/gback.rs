use std::{
  env,
  process::{exit, Command},
};

const RED: &str = "\x1b[31m";
const NORMAL: &str = "\x1b[0m";

fn main() {
  let commit = match env::args().nth(1) {
    Some(v) if !v.is_empty() => v,
    _ => {
      eprintln!("{}❌ Commit hash required{}", RED, NORMAL);
      exit(1);
    }
  };

  let status = Command::new("git")
    .args(["restore", "--source", &commit, "--", "."])
    .status()
    .expect("Failed to run git restore");

  exit(status.code().unwrap_or(1));
}

