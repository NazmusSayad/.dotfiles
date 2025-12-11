use std::{
  io::{self, Write},
  process::{exit, Command},
};

fn main() {
  println!("Restore and clean?");
  print!("Press [Enter] to confirm, or any other key to cancel: ");
  io::stdout().flush().ok();

  let mut input = String::new();
  io::stdin().read_line(&mut input).ok();
  let confirm = input.trim_end_matches(|c| c == '\n' || c == '\r');

  if !confirm.is_empty() {
    eprintln!("❌ Aborted.");
    return;
  }

  let status = Command::new("git")
    .args(["restore", "."])
    .status()
    .expect("Failed to run git restore");

  if !status.success() {
    exit(status.code().unwrap_or(1));
  }

  let status = Command::new("git")
    .args(["clean", "-fd"])
    .status()
    .expect("Failed to run git clean");

  exit(status.code().unwrap_or(1));
}

