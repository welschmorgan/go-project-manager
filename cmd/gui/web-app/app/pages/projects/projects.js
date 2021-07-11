class Projects {
  constructor() {
    this.columns = ['Name', 'Url', 'Path', 'Type', 'SourceControl'];
  }

  init() {
    this.view = document.querySelector('#projects-list-body');

    fetch("/api/projects")
      .then(res => res.json())
      .then(projects => {
        this.view.innerHTML = '';
        for (const project of projects) {
          const row = document.createElement('tr');
          for (const col of this.columns) {
            const cell = document.createElement('td');
            cell.textContent = project[col];
            row.appendChild(cell);
          }
          this.view.appendChild(row);
        }
      })
  }
}