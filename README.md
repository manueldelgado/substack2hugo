![Built with Go](https://img.shields.io/badge/Built%20with-Go-00ADD8?logo=go&logoColor=white)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

# Substack to Hugo Converter

This Go program helps writers and creators who want to **migrate away from [Substack](https://substack.com)** and take full ownership of their online presence by moving to [Hugo](https://gohugo.io), a fast and flexible static site generator.

If you've felt limited by Substack and want **total creative freedom**, **full control of your content**, and **better customization**, this tool is for you.

---

## Why this tool exists

When exporting your data from Substack, the provided format is, in my opinion, **suboptimal**:

- Substack exports a single posts.csv file containing only a subset of the blog post metadata (e.g. - no SEO description).
- The actual content of each post is exported separately as individual HTML files.
- The metadata and the content are not combined into a single file, making it cumbersome to rebuild your site elsewhere.

This program solves that.

It uses the OpenAI Batch API to generate the missing metadata (SEO title, SEO description and SEO keywords), and automatically combines the metadata and content into **Hugo-ready HTML files**, complete with proper front matter, so you can effortlessly recreate your blog on Hugo.


---

## What's New in v1.1.0

- This version enhances functionality by generating **SEO-optimized titles, descriptions, and keywords** for each blog post using **OpenAI GPT-4.1**.
- The process is now divided into **three steps**:

> [!IMPORTANT]
> If you are not interested in generating descriptions or keywords and the existing title works just fine for you (or are unable/unwilling to use the OpenAI API), you might be more interested in the [v1.0.0 version](https://github.com/manueldelgado/substack2hugo/releases/tag/v1.0.0), which does not bother with those enhancements. 

### 1. Create the Batch File

Use the `generate-batch` module to create a `posts2upload.jsonl` file. This file contains the prompts (the original Substack HTML posts combined with a prompt template) to be processed by the OpenAI API.

### 2. Upload and Process the Batch File

I'm not including detailed instructions on how to do this, because the explanation on the [OpenAI Batch API docs](https://platform.openai.com/docs/guides/batch) is great. You can use the code samples they provide and you will only need to change the name of the .jsonl file to upload, and the IDs of the input file, the batch job, and the output file. Also, see below my comments on this API and what you need before you can run this program.

Anyway, the basic steps are:

- Upload `posts2upload.jsonl` from your local folder
- Create a batch job that references the uploaded file
- Wait for OpenAI to process the batch (can take up to 24 hours). Check the docs on how to check the job status.
- Download the output JSONL file, that must be named `batch_output.jsonl` (as in the OpenAI code samples, for simplicity).

### 3. Generate Hugo-Formatted Files

Run the `substack2hugo` program to generate the Hugo-ready HTML files into the `/hugohtml` folder.

You can customize:
- **Draft inclusion**: Use `--ignore-drafts` to skip unpublished posts.
- **Title source**: Use `--use-seo-title` to replace the original Substack title with the SEO-optimized title generated with the OpenAI API.

---

## How It Works

- **Input (provided by you, exported from Substack)**:
  - `posts.csv` 
  - `/posts/` — a folder containing all post HTML files (named with their `post_id`).

- **Input (included here)**:
  - `prompt.txt` - this is the prompt I have used, adapted from polepole's answer to [this question](https://community.openai.com/t/keywords-for-my-article-text/932201) on the OpenAI forums. You might wish to modify the prompt based on your preferences.

- **Output**:
  - A `/hugohtml/` folder will be created.
  - For each post, a new HTML file is generated, containing:
    - Hugo-compatible front matter (metadata).
    - The full original or SEO-enhanced HTML content.

- **Important**:
  - If `/hugohtml` already exists, **all files inside it will be deleted** each time the program runs.
  - This ensures a clean export every time.

- **Drafts Handling**:
  - By default, **all posts** (published and drafts) are processed.
  - Use the `--ignore-drafts` flag to **skip** drafts.

- **Title Handling**:
  - By default, the **original Substack title** is used.
  - Use the `--use-seo-title` flag to **use the SEO-optimized title** instead.

---

## Using the OpenAI Batch API

This project uses the [OpenAI Batch API](https://platform.openai.com/docs/guides/batch) for cost-effective processing.

- **Requirements**:
  - A valid OpenAI platform account.
  - A valid API key.
  - Sufficient credit or balance.

- **Why Batch API?**
  - Batch processing is **50% cheaper** compared to synchronous API calls.
  - But it can take up to **24 hours** to complete processing. My blog -200 articles- only took like one hour, though, for less than $50 cents.

- **Important**: Costs will vary depending on the total size of your content but, for most personal blogs, are usually negligible.

- **Important (2)**: Double check the output from OpenAI. You know how these models work: they sometimes ignore your instructions, become "too creative" or act inconsistently between calls. 

- **Flexibility**: This is the workflow that suits my needs best. Others might prefer:
  - Using synchronous API calls.
  - Using different OpenAI LLMs or another provider.

Feel free to adapt the approach!

---

## Getting Started

1. Install [Go](https://go.dev/dl/) if you haven't already.
2. Clone this repository or download the source code.
3. Make sure your `posts.csv` and `/posts/` folder (from your Substack export) are ready.
4. Follow the three-step process described above.

---

## License and Contributions

This code is licensed under the [MIT License](LICENSE).

Feel free to open an issue if you have feedback or ideas for improvement.
I cannot, however, commit to solving all your doubts or evolve the code, so pull requests are **very welcome**!

---

**Created with the intention of helping creators own their content and build their sites their way.**

---

**Links:**
- [Hugo Static Site Generator](https://gohugo.io)
- [Substack](https://substack.com)
- [OpenAI Batch API Documentation](https://platform.openai.com/docs/guides/batch)
