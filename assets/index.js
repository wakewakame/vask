// タスクの作成
const newTask = (name, expect = "0h-0h", actual = "0h") => {
  return {
    "name": name,      // タスク名
    "child": [],       // 子タスク
    "expect": expect,  // タスクの見積値
    "actual": actual   // タスクの実績値
  };
};

// markdown 形式のタスクリストを json に変換
const parseTask = (text) => {
  const root_task = newTask("root");
  const task_stack = [root_task];
  text.split("\n").forEach((line) => {
    // 行のパース
    const re = /^([\t ]*)-\s*([0-9]+[hd])\s*-\s*([0-9]+[hd])\s*,\s*([0-9]+[hd])\s*:\s*(.*)$/;
    const match = line.match(re);
    if (match === null) { return; }
    const indent = match[1].length;
    const name = match[5];
    const expect = `${match[2]}-${match[3]}`;
    const actual = match[4];

    // インデントが減っていれば task_stack を減らす
    while (indent + 1 < task_stack.length) {
      task_stack.pop();
    }

    // タスクの作成
    const task = newTask(name, expect, actual);
    task_stack[indent].child.push(task)

    // task_stack に追加
    task_stack.push(task);
  });
  return root_task;
};

const getProject = async () => {
  const res = await fetch("./mock/project.md");
  const text = await res.text();
  return parseTask(text);
};

const main = async () => {
  const project = await getProject();
  console.log(JSON.stringify(project, null, 2));
};

main();
