# How to Push to GitHub

Since the GitHub CLI (`gh`) is not installed on this machine, you need to create the repository manually.

1.  **Go to GitHub**: Open [https://github.com/new](https://github.com/new) in your browser.
2.  **Create Repository**:
    *   Repository name: `urumi-store-orchestrator`
    *   Description: "Store Provisioning Platform for Urumi AI Internship"
    *   Select **Public** or **Private**.
    *   **Do NOT** initialize with README, .gitignore, or License (we already have them).
    *   Click **Create repository**.

3.  **Push Code**:
    Run the following commands in your terminal (inside `urumi-store-orchestrator` folder):

    ```powershell
    git remote add origin https://github.com/bansal.bhhunesh/urumi-store-orchestrator.git
    git branch -M main
    git push -u origin main
    ```

    *(Note: Replace `bansal.bhhunesh` with your actual GitHub username if different).*
