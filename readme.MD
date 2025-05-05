# Local Tree 
**Simple CLI tool for generating trees for your ideas or project structure.**

With the rise of MCP and agentic AI I have been exposed to an increasing number of project file structure trees. As the result of this, I have found myself imagining more things in this format and found a want for a tool like this. Sorry the instructions are written like I'm a mac user, I'm a mac user. 

## Usage 
1. Clone the repository `notongithub.comyetlol`
2. build the project 
    - `go build -o localtree`
3. Invoke with `localtree` optionally, provide a file as an argument. 

## Add to PATH 

To use `localtree` from anywhere in your terminal, add its directory to your PATH.  
If you built the binary in your project folder, run:

```sh
echo 'export PATH="$PATH:$(pwd)"' >> ~/.zshrc
source ~/.zshrc
```
- For Bash, use `~/.bashrc` instead of` ~/.zshrc`
- For Fish shell, use:
    - `set -U fish_user_paths (pwd) $fish_user_paths`

## Features and specs 
- Implemented in GO
- TUI experience with file I/O
- Simple to use, like scratch paper for trees

### Example Tree 
```
Root
├──Notebooks
├─── Granite4.0.ipynb
│  └─── Transformers notebook for running tiny granite 4.0.
└─── DirectoryGuide.md
   └─── Markdown tree of the folder (This file!)
```

Inspired entirley by https://tree.nathanfriend.com. Wanted a local, TUI tool that served the same purpose. Lots of this project was generated with GPT 4.1 using copilot in vscode with mcp. I *do not* take full credit for the code in this repository. 
