![Built with Go](https://img.shields.io/badge/Built%20with-Go-00ADD8?logo=go&logoColor=white)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

# Substack to Hugo Converter

This Go program helps writers and creators who want to **migrate away from** [**Substack**](https://substack.com) and take full ownership of their online presence by moving to [Hugo](https://gohugo.io), a fast and flexible static site generator.

If you've felt limited by Substack and want **total creative freedom**, **full control of your content**, and **better customization**, this tool is for you.

---

## Why This Tool Exists

When exporting your data from Substack, the provided format is, in my opinion, **suboptimal**:

- Substack exports a single `posts.csv` file containing only a subset of the blog post metadata (e.g. - no SEO description).
- The actual content of each post is exported separately as individual HTML files.
- The metadata and the content are **not combined** into a single file, making it cumbersome to rebuild your site elsewhere.

This program solves that.

It automatically combines the metadata and content into **Hugo-ready HTML files**, complete with proper front matter, so you can effortlessly recreate your blog on Hugo.

---

## How It Works

- **Input (Provided by You)**:

  - `posts.csv` — exported from Substack.
  - `/posts/` — a folder containing all post HTML files (named with their `post_id`).

- **Output**:

  - A `/hugohtml/` folder will be created.
  - For each post, a new HTML file is generated, containing:
    - Hugo-compatible front matter (metadata).
    - Two blank lines.
    - The full original HTML content.

- **Important**:

  - If `/hugohtml` already exists, **all files inside it will be deleted** each time the program runs.
  - This ensures a clean export every time.

- **Drafts Handling**:

  - By default, **all posts** (published and drafts) are processed.
  - You can pass the `--ignore-drafts` flag when running the program to **skip posts** that were unpublished on Substack.

Example usage:

```bash
go run substack2hugo.go                      # Process all posts (published + drafts)
go run substack2hugo.go --ignore-drafts       # Process only published posts
```

The program outputs a summary listing:

- Number of successfully created files.
- Any errors encountered.
- The absolute path to the `/hugohtml` folder.

---

## Running the Program

⚡ **Important:**\
You must run the program **inside the same folder** where your Substack export files (`posts.csv` and `/posts/`) are located.

For example:

```bash
cd /path/to/your/substack-export-folder
go run /path/to/substack-to-hugo/substack2hugo.go
```

Or compile it first and then run:

```bash
cd /path/to/your/substack-export-folder
/path/to/substack-to-hugo/substack2hugo --ignore-drafts
```

---

## Getting Started

1. Install [Go](https://go.dev/dl/) if you haven't already.
2. Clone this repository or download the source code.
3. Make sure your `posts.csv` and `/posts/` folder (from Substack export) are ready.
4. Run the program from that directory as shown above.

---

## License and Contributions

This code is licensed under the [MIT License](LICENSE).

Enhancements, suggestions, and pull requests are **very welcome**! Feel free to open an issue if you have feedback or ideas for improvement.

---

**Created with the intention of helping creators own their content and build their sites their way.**

---

**Links:**

- [Hugo Static Site Generator](https://gohugo.io)
- [Substack](https://substack.com)

