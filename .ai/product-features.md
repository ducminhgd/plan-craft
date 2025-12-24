# Product features

This project has some core features:

1. Project
2. Work Breakdown Structure (WBS)
3. Timeline & Dependencies
4. Resources
5. Costs

## MVP Features (Must-Have)

These define the minimum usable product.

1. Project Definition
   1. Project metadata (name, type, methodology: Waterfall / Agile / Hybrid)
   2. Start date, target end date
   3. Assumptions & constraints
2. Task & WBS Management
   1. Hierarchical tasks (epics → tasks → subtasks)
   2. Estimated effort (hours / days)
   3. Dependency types:
      1. Finish-to-Start
      2. Start-to-Start
   4. Critical path calculation (basic)
3. Timeline Estimation
   1. Auto-calculated project duration
   2. Gantt-style timeline (read-only initially)
   3. Slack / buffer visibility
4. Resource Planning
   1. Define roles (PM, Backend, QA, Designer, etc.)
   2. Assign resources to tasks
   3. Capacity limits per resource
5. Cost Estimation
   1. Hourly / daily cost per role
   2. Total cost per task, phase, project
   3. Cost breakdown by category

## What end-users can do with Plan Craft

1. Input work items for a projects, can estimate how many man-day or man-month for each work item, and for each milestone, and for each project.
2. Can estimate the timeline or roadmap of the project. Because there are some work items depend on others.
3. From estimated man-day or man-month, and the timeline, we can estimate how much people we need to work on the project, and for how long.
4. From the number of people we need, and the cost per person, we can estimate the cost of the project.

## Roadmap

1. Version 1.0: Project management, work items management.
2. Version 1.1: Timeline estimation.
3. Version 1.2: Resource planning.
4. Version 1.3: Cost estimation.