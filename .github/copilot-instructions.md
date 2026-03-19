Keep the code simple.
Focus on readability and maintainability. 
Always make code readable even if it becomes longer
Avoid adding comments to code unless explicitly requested.
Do not remove comments that are already present created by the user.
Instead of using "err" try to add more context for readability "userTokenErr"
Be explicit when choosing names. 

Here is a link to the go templ template llm instructions: https://templ.guide/llms.md


Check that you did not wright clever code. 

CI notes:
- GitHub Actions runs lint and build/test on pull requests and on pushes to main.
- Deploy runs only on push to main (via buildAndTest workflow).
- Workflows generate templ files (templ v0.3.960) before lint/build; keep templ-generated files in mind when troubleshooting CI typecheck issues.
