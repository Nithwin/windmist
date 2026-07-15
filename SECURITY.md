# Security Policy for WindMist 🌪️

We take the security and integrity of **WindMist** seriously. Because WindMist is an autonomous terminal AI coding assistant—with capabilities to read/write files, execute system commands, and communicate with remote AI models—maintaining robust security boundaries is a top priority.

---

## 🛡️ Supported Versions

We actively maintain and provide security patches for the following versions of WindMist:

| Version | Supported          | Remarks |
| :---    | :---               | :--- |
| `main` (latest) | :white_check_mark: Yes | Active development branch; security patches applied immediately. |
| `1.x.x` | :white_check_mark: Yes | Current stable release series. |
| `< 1.0.0` (Dev/Pre-release) | :x: No | Early development snapshots; upgrade to `main` or latest release. |

---

## 🚨 Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub Issues, Discussions, or Pull Requests.**

If you discover a potential security vulnerability, please report it privately using one of the following methods:

### 1. GitHub Private Vulnerability Reporting (Preferred)
1. Navigate to the **[Security tab](https://github.com/your-username/windmist/security)** of this repository.
2. Click on **Advisories** in the left sidebar.
3. Click **Report a vulnerability** to open a confidential security advisory.

### 2. Direct Contact
If you are unable to use GitHub Security Advisories, please contact the repository owner directly through GitHub or open a private security/moderation report. Please include `[SECURITY]` in the subject line or message prefix.

### What to Include in Your Report
To help us triage and resolve the issue quickly, please provide:
- **Type of Vulnerability:** (e.g., Path Traversal, Arbitrary Command Execution, API Key Leakage).
- **Affected Layer/Component:** Specify whether it occurs in the agent loop, tool execution, provider integration, or terminal UI.
- **Steps to Reproduce:** A step-by-step reproduction guide or minimal Proof-of-Concept (PoC) script/prompt.
- **Impact Assessment:** How an attacker could exploit this bug and what system access or data would be compromised.
- **Environment Details:** Go version (`go version`) and operating system (Linux/macOS/Windows).

---

## ⏳ Response & Disclosure Timeline

We are committed to being responsive and transparent with security researchers:

- **Initial Acknowledgement:** Within **48 hours** of receiving the report.
- **Triage & Assessment:** Within **7 days**, confirming whether the issue is verified and establishing its severity rating (CVSS).
- **Patch Development:** A fix or mitigation will be developed and tested in a private security branch.
- **Coordinated Disclosure:** Once a patch is released, we will publish a public GitHub Security Advisory crediting the reporter (unless anonymity is requested).

---

## 🔒 Key Threat Models & Security Considerations

Because WindMist is an autonomous terminal AI coding assistant with direct access to filesystem operations, command execution, and remote model providers, we pay special attention to the following attack vectors:

### 1. Workspace Boundaries (`Workspace boundaries`)
File operations (`create`, `read`, `write`, `edit`) must strictly confine themselves to the target project workspace. Vulnerabilities allowing relative path escape (e.g., `../../etc/passwd` or overwriting `~/.bashrc` without explicit user override) are classified as **High/Critical severity**.

### 2. Tool Permissions (`Tool permissions`)
Tools that interact with the system or filesystem must operate strictly within their granted capabilities and validation rules. Any flaw allowing a tool to exceed its intended scope, manipulate unauthorized files, or bypass access boundaries is treated as a high-priority vulnerability.

### 3. Unsafe Commands (`Unsafe commands`)
WindMist can execute terminal commands on behalf of the user. Commands flagged as potentially destructive (`SafeToAutoRun = false`) **must** require explicit user approval via the terminal UI before executing. Any vulnerability that allows untrusted model output or external prompts to bypass user confirmation for unsafe commands is classified as a **Critical vulnerability**.

### 4. API Keys & Secret Protection (`API keys`)
WindMist handles sensitive API keys for external providers (e.g., `GEMINI_API_KEY`, `OPENAI_API_KEY`). API keys must **never** be logged to stdout/stderr, written to disk in plain text logs, or exposed in debug traces. HTTP authorization headers must be stripped or masked during logging.

### 5. Prompt Injection (`Prompt injection`)
While Large Language Models (LLMs) are inherently susceptible to indirect prompt injection from untrusted repository files, WindMist's deterministic engine acts as a non-bypassable guardrail. Even if a model is manipulated via prompt injection, the engine must enforce strict boundary checks and user confirmation workflows before executing high-risk system actions.

---

## ⚖️ Responsible Disclosure & Safe Harbor

If you conduct your security research and vulnerability reporting in accordance with this policy, we consider your actions to be authorized and will not initiate legal action against you. We ask that you:
- Make a good-faith effort to avoid privacy violations, destruction of user data, or degradation of services during testing.
- Only interact with test repositories and environments you own or have explicit permission to test against.
- Do not exploit a security vulnerability beyond what is strictly necessary to demonstrate the flaw.

Thank you for helping keep **WindMist** and its community safe! 🚀🌪️
