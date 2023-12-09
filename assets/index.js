"use strict";

const Task = class {
  // 以下のような markdown をパースする
  // - 3h-4.5h, 0.5h: task
  //     - 1h-2h, 0.5h: task1
  //     - 2h-3.5h, 0h: task2
  static parseMarkdown(markdown) {
    const rootTask = new Task("root");
    const taskStack = [{task: rootTask, indent: -1}];
    const errs = [];
    const re = /^([\t ]*)-\s*([0-9]+(\.[0-9]+)?[mh])\s*-\s*([0-9]+(\.[0-9]+)?[mh])\s*,\s*([0-9]+(\.[0-9]+)?[mh])\s*:\s*(.*)$/;
    markdown.split("\n").forEach((line, lineNumber) => {
      // 行のパース
      const match = line.match(re);
      if (match === null) {
        const err = `failed to parse at L${lineNumber + 1}: ${line}`;
        errs.push(err);
        return;
      }
      const indent = match[1].length;
      const expect = `${match[2]}-${match[4]}`;
      const actual = match[6];
      const name = match[8];
      const task = new Task(name, expect, actual);

      // タスクの追加
      while (indent <= taskStack.slice(-1)[0].indent) {
        taskStack.pop();
      }
      taskStack.slice(-1)[0].task.child.push(task)
      taskStack.push({task, indent});
    });
    return [rootTask, errs];
  }

  constructor(name, expect = "0h-0h", actual = "0h") {
    this.name = name;
    this.child = [];
    this.expect = expect;
    this.actual = actual;
  }

  toMarkdown() {
    const format = (task, indent) => {
      const indentStr = [...Array(indent)].map(() => "    ").join("");
      return [
        `${indentStr}- ${task.expect}, ${task.actual}: ${task.name}`,
        ...task.child.map(t => format(t, indent + 1))
      ].join("\n");
    };
    return [...this.child.map(task => format(task, 0))].join("\n");
  }

  // JSON.stringify() 時に呼び出されるメソッド
  toJSON() {
    return {
      name: this.name,
      child: this.child,
      expect: this.expect,
      actual: this.actual
    };
  }
};

const getProject = async () => {
  const res = await fetch("./mock/project.md");
  const text = await res.text();
  return Task.parseMarkdown(text);
};

const main = async () => {
  const [project, _errs] = await getProject();
  console.log(JSON.stringify(project, null, 2));
  console.log(project.toMarkdown());
};

main();
