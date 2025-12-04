---
name: Bug report
title: "[Bug]"
description: Create a new bug report
projects: []
labels:
  - bug
body:
  - type: textarea
    id: text-steps-to-reproduce
    attributes:
      label: Steps to Reproduce
      description: Describe the steps you took to create the bug. Add screenshots to
        help explain your problem.
    validations:
      required: true
  - type: textarea
    id: text-expected-behavior
    attributes:
      label: Expected Behaviour
      description: Describe what you expected to happen in this scenario.
    validations:
      required: true
  - type: textarea
    id: text-additional-information
    attributes:
      label: Additional Information
      description: "Provide any other context to the problem. Device, browser,
        localization, and similar. "
---
