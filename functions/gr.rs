use std::{
  io::{self, Read, Write},
  process::{exit, Command},
};

const RED: &str = "\x1b[31m";
const DIM: &str = "\x1b[2m";
const NORMAL: &str = "\x1b[0m";

fn main() {
  println!("{}Restore and clean?{}", RED, NORMAL);
  print!("{}Press [Enter] to confirm, or any other key to cancel: {}", DIM, NORMAL);
  io::stdout().flush().ok();

  let mut buf = [0u8; 1];
  let read = io::stdin().read(&mut buf).unwrap_or(0);
  let confirm = read == 0 || buf[0] == b'\n' || buf[0] == b'\r';

  if !confirm {
    eprintln!("{}❌ Aborted.{}", RED, NORMAL);
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

