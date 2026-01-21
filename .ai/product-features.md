# Product features

This product has some core features:

1. Project management
   1. Project metadata (name, type)
   2. Start date, target end date
   4. Client (we have client management)
   5. Configurations
      1. Hours Per Day
      2. Days Per Week
      3. Working days per week (Monday to Friday by default)
      4. Timezone
      5. Currency
2. Work Breakdown Structure (WBS)
   1. A milestone has:
      1. Name
      2. Project it belongs to (required)
      3. Start date
      4. End date
      5. Status
      6. Description
   2. A task has:
      1. Name
      2. Level: if the task does not belong to a parent task, the level is 1. If the task belongs to a parent task, the level is the parent task's level + 1.
      3. Project it belongs to (required)
      4. Milestones it belongs to (optional). When a task is added to a milestone, it will be added to the project of that milestone automatically.
      5. Priority: Low, Medium, High, Critical.
      6. Estimated effort in man-days
      7. Parent task (optional). A task belongs to a parent task only.
      8. Description
      9. Status: To Do, In Progress, Done, Cancelled.
   3. A task can be linked to another task with a dependency. The dependency can be:
      1. Blocking: a task cannot be started until this task is completed.
      2. Blocked by: this task cannot be started until another task is completed.
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
   2. A resource can be assigned to multiple projects. For each projects, a resource can be assigned to multiple roles.
   3. A resource can be allocated in a project with a spefic percentage for a specific role in a specific time range.
5. High-level estimation. We don't need to assign tasks to a specific resource, we just need to estimate how many man-days or man-months for a task, or a projects.
   1. The estimation can be summed up to the milestone, and the project level.
6. Project costs: include costs for human resources and cost for infrastructure & service costs.
   1. Human resources
      1. Hourly / daily cost per role
      2. Total cost per task, phase, project
   2. Cost breakdown by category
   3. The cost can be summed up to the milestone, and the project level.
   4. The costs can be filtered by categories, by resources, by milestones, by dates.
7. The menu bar:
   1. `File`:

      ```
      File
      ├── Open file (Ctrl/Cmd+O)
      ├── Save as (Ctrl/Cmd+Shift+S)
      ├── ------------------------ (seperation line)
      └── Exit (Ctrl/Cmd+Q)
      ```

      1. `File` > `Open file`: to open a project SQLite3 database file. The application will close the current project if there is any, and open the new project. By default, we don't point to any database file. When a database is opened, we will auto-save the project to that database file.
      2. `File` > `Save as`: to save the current project to a new SQLite3 database file.
      3. `File` > `Exit`: to close the application.
   2. `Help`:

      ```
      Help
      ├── Guides
      ├── ------------------------ (seperation line)
      └── About
      ```

      1. `Help` > `Guides`: to open the guides page in the browser.
      2. `Help` > `About`: to show the about dialog.

## What end-users can do with Plan Craft

1. Input work items for a projects, can estimate how many man-day or man-month for each work item, and for each milestone, and for each project.
2. Can estimate the timeline or roadmap of the project. Because there are some work items depend on others.
3. From estimated man-day or man-month, and the timeline, we can estimate how much people we need to work on the project, and for how long.
4. From the number of people we need, and the cost per person, we can estimate the cost of the project.

## Project plan

### Project Roles

1. A project has many roles. Each role has a name, and a level (junior, mid, senior, lead), and headcount.
2. The roles are managed in the Project Form, there is a table called `Roles`.
3. Unique on name and level, that means each role of each level can only have 1 row.
4. In the `Project` > `Roles` page, we can add, edit, and delete roles.

### Human Resource planning

1. For each Task, we can assign roles and estimate the effort in man-days. The roles are not limited.
2. In the Form view of Task, there is a section called `Resources`. It's a table with 3 columns: `#`, `Role`, `Effort`, `Action`.
   1. `#` is the index of the role. It's auto-incremented when adding a new role.
   2. `Role` is the name of the role. Selected from a list of roles from the project.
   3. `Effort` is the effort in man-hours for this role. It's a number input field.
   4. `Action` is a dropdown with 2 options: `Remove`, and `Edit`.
   5. There is a button `Add row` to add a new role to the task.
   6. Unique by roles.
   7. In the header, there are total roles, total effort in man-hours, and total effort in man-days.
