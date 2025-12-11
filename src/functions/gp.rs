use std::{
  env,
  process::{exit, Command},
};

fn main() {
  let args: Vec<String> = env::args().skip(1).collect();

  let current_branch = Command::new("git")
    .args(["branch", "--show-current"])
    .output()
    .map(|o| String::from_utf8_lossy(&o.stdout).trim().to_string())
    .unwrap_or_default();

  let target_branch = match args.len() {
    0 => {
      println!("No branch specified, using current branch");
      current_branch.clone()
    }
    1 => args[0].clone(),
    _ => {
      eprintln!("Usage: gp [branch]");
      exit(1);
    }
  };

  println!(
    "Pulling changes from {} into {} (default)",
    target_branch, current_branch
  );

  let status = Command::new("git")
    .args(["prune", "--progress"])
    .status()
    .expect("Failed to run git prune");

  if !status.success() {
    exit(status.code().unwrap_or(1));
  }

  let status = Command::new("git")
    .args(["pull", "origin", &target_branch, "--progress"])
    .status()
    .expect("Failed to run git pull");

  exit(status.code().unwrap_or(1));
}

