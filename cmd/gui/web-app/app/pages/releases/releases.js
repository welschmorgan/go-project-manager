function formatDate(d) {
  d = typeof(d) === 'string' ? new Date(d) : d;
  return d.getDate() + '/' + (d.getMonth() + 1) + '/' + d.getFullYear()
    + ' '
    + d.getHours() + ':' + d.getMinutes();
}

class Releases {
  constructor() {
    this.releaseTableColumns = [
      {getter: (obj) => obj.project.name},
      {getter: (obj) => obj.project.type},
      {getter: (obj) => formatDate(obj.context.date)},
      {getter: (obj) => obj.context.state == 63 ? 'finished' : 'unfinished'},
      {getter: (obj) => obj.project.source_control},
    ];
    this.undoActionsTableColumns = [
      {getter: (obj) => obj.id},
      {getter: (obj) => formatDate(obj.date)},
      {getter: (obj) => obj.name},
      {getter: (obj) => obj.title},
      {getter: (obj) => obj.params},
    ];
  }

  init() {
    this.releaseTable = document.querySelector('#release-table-body');
    this.undoActionsTable = document.querySelector('#undo-actions-table-body');
    this.releaseVersionList = document.querySelector('#version-list');
    this.debug = document.querySelector('#debug');
    fetch("/api/release/")
      .then(async res => {
        const text = await res.text();
        this.debug.innerText = text;
        return JSON.parse(text);
      })
      .then(releases => {
        this.releases = releases;
        this.renderVersionList();
        this.renderRelease(Object.keys(this.releases)[0]);
      })
      .catch((reason) => console.error(reason));
  }

  renderVersionList() {
    const createItem = (value) => {
      const item = document.createElement('li');
      const anchor = document.createElement('a');
      anchor.href = '#';
      anchor.innerText = value;
      item.addEventListener('click', () => this.renderRelease(value), true);
      item.appendChild(anchor);
      this.releaseVersionList.appendChild(item);
    };
    this.releaseVersionList.innerHTML = '';
    for (const version in this.releases) {
      createItem(version);
    }
  }

  renderRelease(version) {
    const createCell = (row, value) => {
      const cell = document.createElement('td');
      cell.innerHTML = value;
      row.appendChild(cell);
      return cell;
    };
    this.releaseTable.innerHTML = '';
    const release = this.releases[version];
    for (const action of release) {
      const row = document.createElement('tr');
      for (const col of this.releaseTableColumns) {
        createCell(row, col.getter(action));
      }
      this.releaseTable.appendChild(row);
    }
    this.renderUndoActions(version);
  }

  renderUndoActions(version) {
    const createCell = (row, value) => {
      const cell = document.createElement('td');
      cell.innerHTML = value;
      row.appendChild(cell);
      return cell;
    };
    this.undoActionsTable.innerHTML = '';
    const releases = this.releases[version];
    for (const release of releases) {
      const actions = release.undoActions;
      for (const action of actions) {
        const row = document.createElement('tr');
        for (const col of this.undoActionsTableColumns) {
          createCell(row, col.getter(action));
        }
        this.undoActionsTable.appendChild(row);
      }
    }
  }
}