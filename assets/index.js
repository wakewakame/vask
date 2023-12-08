const Project = class {
  constructor(markdown, title) {
    this.rootTask = Project.parseMarkdown(markdown, title);
  }

  toJSON() {
    return this.rootTask;
  }

  toMarkdown() {
    return this.rootTask.toMarkdown();
  }

  static parseMarkdown(markdown, title) {
    const rootTask = new Task(title ?? "root");
    const taskStack = [rootTask];
    const re = /^([\t]*)-\s*([0-9]+(\.[0-9]+)?[mh])\s*-\s*([0-9]+(\.[0-9]+)?[mh])\s*,\s*([0-9]+(\.[0-9]+)?[mh])\s*:\s*(.*)$/;
    let lineNumber = 0;
    for (let line of markdown.split("\n")) {
      lineNumber += 1;

      // 行のパース
      if (line === "") { continue; }
      const match = line.match(re);
      if (match === null) {
        console.error(`failed to parse at L${lineNumber}: ${line}`);
        continue;
      }
      const indent = match[1].length;
      const [name, expect, actual] = [match[8], `${match[2]}-${match[4]}`, match[6]];

      // インデントが減っていれば taskStack を減らす
      while (indent + 1 < taskStack.length) {
        taskStack.pop();
      }

      // タスクの作成
      const task = new Task(name, expect, actual);
      taskStack[indent].child.push(task)

      // taskStack に追加
      taskStack.push(task);
    }
    return rootTask;
  }
};

// タスクの作成
const Task = class {
  constructor(name, expect = "0h-0h", actual = "0h") {
    this.name = name;
    this.child = [];
    this.expect = expect;
    this.actual = actual;
  }

  toJSON() {
    return {
      "name": this.name,
      "child": this.child,
      "expect": this.expect,
      "actual": this.actual
    };
  }

  toMarkdown() {
    return [
      `- ${this.expect}, ${this.actual}: ${this.name}`,
      ...this.child.map(task =>
        task.toMarkdown().split("\n").map(line =>
          line.replace(/^/, " ")
        ).join("\n")
      )
    ].join("\n");
  }
};

const getProject = async () => {
  const res = await fetch("./mock/project.md");
  const text = await res.text();
  return new Project(text, "tmp");
};

const main = async () => {
  const project = await getProject();
  console.log(JSON.stringify(project, null, 2));
  console.log(project.toMarkdown());
};

main();
