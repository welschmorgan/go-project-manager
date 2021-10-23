class Versions {
  constructor() {
    this.columns = ['Name', 'Version']
  }

  init() {
    this.view = document.querySelector('#versions-list-body');

    fetch("/api/versions")
      .then(res => res.json())
      .then(versions => {
        this.view.innerHTML = '';
        for (const project of versions) {
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