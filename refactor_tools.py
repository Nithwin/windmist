import os
import glob
import re

for file_path in glob.glob('/home/shadow/Desktop/windmist/internal/tools/**/*.go', recursive=True):
    with open(file_path, 'r') as f:
        content = f.read()

    if 'func (' not in content or 'Definition() tools.Definition' not in content:
        continue

    category = "tools.CategoryFilesystem"
    perm = "tools.PermReadOnly"

    if "/editing/" in file_path:
        category = "tools.CategoryEditing"
        if "search" in file_path:
            category = "tools.CategorySearch"
            perm = "tools.PermReadOnly"
        else:
            perm = "tools.PermWrite"
    elif "/filesystem/" in file_path:
        if "glob" in file_path or "grep" in file_path:
            category = "tools.CategorySearch"
        elif "delete" in file_path or "write" in file_path or "append" in file_path or "create" in file_path or "rename" in file_path:
            perm = "tools.PermWrite"
    elif "/system/" in file_path:
        if "git" in file_path:
            category = "tools.CategoryGit"
            perm = "tools.PermDangerous"
        else:
            category = "tools.CategorySystem"
            perm = "tools.PermDangerous"
    elif "/web/" in file_path:
        category = "tools.CategoryWeb"
    elif "/agent/" in file_path:
        category = "tools.CategoryAgent"
        perm = "tools.PermWrite"

    # Match `tools.Definition{\n\t\tName: "...",\n\t\tDescription: "...",`
    
    def repl(m):
        return f"{m.group(0)}\n\t\tCategory:   {category},\n\t\tPermission: {perm},"
        
    new_content = re.sub(r'(tools\.Definition\{\s*Name:\s*".*?",\s*Description:\s*".*?",)', repl, content, count=1)
    
    if new_content != content:
        with open(file_path, 'w') as f:
            f.write(new_content)
        print(f"Updated {file_path}")
