use std::{
  env,
  process::{exit, Command, Stdio},
};

const DIM: &str = "\x1b[2m";
const RED: &str = "\x1b[31m";
const NORMAL: &str = "\x1b[0m";

fn main() {
  let args: Vec<String> = env::args().skip(1).collect();
  let repo_path = args.get(0).map(|s| s.as_str()).unwrap_or("");

  let mut parts = repo_path.split('/');
  let first = parts.next().unwrap_or("");
  let second = parts.next();
  let valid = !first.is_empty() && second.map_or(true, |s| !s.is_empty() && parts.next().is_none());

  if valid {
    println!("{}Using GitHub CLI to resolve URL...{}", DIM, NORMAL);

    match Command::new("gh")
      .args(["repo", "view", repo_path, "--json", "url", "-q", ".url"])
      .stderr(Stdio::inherit())
      .output()
    {
      Ok(output) => {
        if output.status.success() {
          let url = String::from_utf8_lossy(&output.stdout).trim().to_string();

          if !url.is_empty() {
            println!("{}GitHub URL: {}{}", DIM, url, NORMAL);

            let mut cmd = Command::new("git");
            cmd.arg("clone").arg(url);

            for arg in args.iter().skip(1) {
              cmd.arg(arg);
            }

            let status = cmd.status().expect("Failed to run git clone");

            exit(status.code().unwrap_or(1));
          }
        }
      }
      Err(err) => {
        eprintln!("{}Failed to run gh: {}{}", RED, err, NORMAL);
      }
    }

    println!(
      "{}Failed to resolve GitHub URL, trying to clone directly...{}",
      RED, NORMAL
    );
  }

  let mut cmd = Command::new("git");
  cmd.arg("clone");

  for arg in args.iter() {
    cmd.arg(arg);
  }

  let status = cmd.status().expect("Failed to run git clone");

  exit(status.code().unwrap_or(1));
}
