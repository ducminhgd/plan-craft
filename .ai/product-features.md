# Product features

This product has some core features:

1. Project management
   1. Project metadata (name, type)
   2. Start date, target end date
   3. Assumptions & constraints
   4. Client (we have client management)
   5. Project Manager
2. Work Breakdown Structure (WBS)
3. Timeline & Dependencies
   1. Auto-calculated project duration
   2. Gantt-style timeline (read-only initially)
   3. Slack / buffer visibility
4. Resources: manage resource, person names, and their titles, and their levels, and their costs.
   1. Resource planning:
      1. Define roles (PM, Backend, QA, Designer, etc.)
      2. Assign resources to tasks
      3. Capacity limits per resource
      4. A role in a project must have the level of the role.
      5. A role in a tasks must be estimated by man-days.
      6. A role in a project must be estimated by man-months, and it can be summed up from its tasks.
5. High-level estimation. We don't need to assign tasks to a specific resource, we just need to estimate how many man-days or man-months for a task, or a projects.
   1. The estimation can be summed up to the milestone, and the project level.
6. Project costs: include costs for human resources and cost for infrastructure & service costs.
   1. Human resources
      1. Hourly / daily cost per role
      2. Total cost per task, phase, project
   2. Cost breakdown by category
   3. The cost can be summed up to the milestone, and the project level.
   4. The costs can be filtered by categories, by resources, by milestones, by dates.

## What end-users can do with Plan Craft

1. Input work items for a projects, can estimate how many man-day or man-month for each work item, and for each milestone, and for each project.
2. Can estimate the timeline or roadmap of the project. Because there are some work items depend on others.
3. From estimated man-day or man-month, and the timeline, we can estimate how much people we need to work on the project, and for how long.
4. From the number of people we need, and the cost per person, we can estimate the cost of the project.

## Roadmap

1. Version 1: Desktop Application can run on MacOS, Windows, and Linux.
   1. Version 1.0: Project management, work items management.
   2. Version 1.1: Timeline estimation.
   3. Version 1.2: Resource planning.
   4. Version 1.3: Cost estimation.
2. Version 2: Web Application. Can connect to an REST API Server. (will do later)

We are implementing version 1 first.