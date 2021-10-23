class Undos {
  constructor() {
    this.columns = [
      {getter: (obj) => obj.name}, 
      {getter: (obj) => obj.title}, 
      {getter: (obj) => obj.path}, 
      {getter: (obj) => {
        const params = Object.entries(obj.params).map(v => `<tr><td>${v[0]}</td><td>${v[1]}</td></tr>`);
        debugger;
        return '<table>' + params.join('') + '</table>';
      }}
    ];
  }

  init() {
    this.tableBody = document.querySelector('#undo-table-body');
    // this.loader = new Loader(this.tableBody);
    fetch("/api/undos")
      .then(res => res.json())
      .then(undos => {
        this.undos = undos;
        this.renderUndos();
      })
      .catch((reason) => console.error(reason));
  }

  renderUndos() {
    const createCell = (row, value) => {
      const cell = document.createElement('td');
      cell.innerHTML = value;
      row.appendChild(cell);
      return cell;
    };
    for (const file in this.undos) {
      const actions = this.undos[file];
      for (const action of actions) {
        const row = document.createElement('tr');
        createCell(row, file);
        for (const col of this.columns) {
          createCell(row, col.getter(action));
        }
        this.tableBody.appendChild(row);
      }
    }
  }
}