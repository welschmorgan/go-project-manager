class Projects {
  constructor() {
    this.columns = ['Name', 'Url', 'Path', 'Type', 'SourceControl'];
    this.projects = [];
  }

  init() {
    this.view = document.querySelector('#projects-list-body');
    this.getProjects()
  }

  renderProjects() {
    for (const project of this.projects) {
      const row = document.createElement('tr');
      for (const col of this.columns) {
        const cell = document.createElement('td');
        cell.textContent = project[col];
        row.appendChild(cell);
      }
      this.view.appendChild(row);
    }
    return this.projects;
  }

  async getProjects() {
    return fetch("/api/projects")
      .then(res => res.json())
      .then(projects => {
        this.view.innerHTML = '';
        this.projects = projects || [];
        this.renderProjects();
      });
  }

  async scanProjects() {
    return fetch('/api/projects/scan', {method: 'POST'})
      .then(res => res.json())
      .then(projects => {
        this.view.innerHTML = '';
        this.projects = projects || [];
        return projects;
      })
      .then(res => this.renderProjects());
  }
}