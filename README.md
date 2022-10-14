# Rapporter

Infrastructure health report generator

This is tool designed to generate reports about low stakes infrastructures that do not have 24/7 support staff (think home lab, or small business).

Instead of constantly probing infrastructure for its state the workflow is:

* The component (let's say backup job) submits its state (few data points + markdown description) to the server.
* The ones that "expire" without being explicitly deleted are marked as failed, to mark components that stopped responding 
* The server compiles a report html about every component, sorted by priority and state.
* The report can be viewed, or sent via e-mail.

The goal is to agument existing monitoring system to cover for rare-but-important tasks that might possibly be annoying to report via normal monitoring.

For example, if you have a cron job this tool could report state of that job and alert if it is running with error or not at all.

