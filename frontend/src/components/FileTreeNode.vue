<script lang="ts" setup>
import type { TreeDropInfo, TreeOption } from "naive-ui";
import { NTree } from "naive-ui";
import { main } from "../../wailsjs/go/models";
import { nextTick, ref, watchEffect } from "vue";
import { Move } from "../../wailsjs/go/main/App.js";
import { promiseResult } from "../utils/promiseResult";

interface TreeOptionWithDir extends TreeOption {
  directory?: string;
}

const props = defineProps<{ dirNode: main.dirNode }>();
const emit = defineEmits(["file-click", "get-snippet-paths"]);

const dataRef = ref<TreeOption[]>([]);
const filesMap = new Map<string, string[]>();

const createData = (
  node: main.dirNode,
  k?: string
): Array<TreeOptionWithDir> => {
  if (!node) return [];

  const files: TreeOptionWithDir[] = node.files.map((file) => {
    let pathParts = file.split("/");
    let fileName: string = "";
    let dir: string = "";
    if (pathParts.length > 1) {
      fileName = pathParts.slice(-1).join("");
      pathParts = pathParts.slice(0, -1);
      dir = pathParts.join("/");
    } else {
      // set to file instead because its in the root path
      fileName = file;
    }

    const fmKey = k != undefined ? `${dir}` : ".";
    if (!filesMap.has(fmKey)) {
      filesMap.set(fmKey, []);
    }
    const filesArr = filesMap.get(fmKey)!;
    filesArr.push(file);

    return {
      directory: dir,
      label: fileName,
      key: file,
      isLeaf: true,
      children: [],
    };
  });

  const subdirs: TreeOptionWithDir[] = Object.entries(node.subdirs || {}).map(
    ([name, subdir]) => {
      // initialize map to avoid skipping empty folder from being listed in map
      if (!filesMap.has(name)) {
        filesMap.set(name, []);
      }
      return {
        label: name,
        key: k != undefined ? `${k}/${name}` : name,
        isLeaf: false,
        children: createData(subdir, name),
      };
    }
  );

  return [...subdirs, ...files];
};

const findSiblingsAndIndex = (
  node: TreeOption,
  nodes?: TreeOption[]
): [TreeOption[], number] | [null, null] => {
  if (!nodes) return [null, null];
  for (let i = 0; i < nodes.length; ++i) {
    const siblingNode = nodes[i];
    if (siblingNode.key === node.key) return [nodes, i];
    const [siblings, index] = findSiblingsAndIndex(node, siblingNode.children);
    if (siblings && index !== null) return [siblings, index];
  }
  return [null, null];
};

const handleDrop = async ({ node, dragNode, dropPosition }: TreeDropInfo) => {
  if (dropPosition == "inside" && !node.isLeaf) {
    // console.log("[inside a folder]", dragNode.label, "to", node.key);
    const [error, _] = await promiseResult(
      Move(dragNode.key!.toString(), node.key!.toString())
    );
    if (error) {
      console.error(error);
      return;
    }
    // console.info(
    //   `successfully moved ${dragNode.key?.toString()} to ${node.key!.toString()}`
    // );
  } else {
    // get the current folder path
    const [nodeSiblings, nodeIndex] = findSiblingsAndIndex(node, dataRef.value);
    if (nodeSiblings === null || nodeIndex === null) {
      return;
    }

    // set to dot as fallback getting the root dir
    let dst: string = ".";
    const pathParts = node.key?.toString().split("/").slice(0, -1);
    if (!pathParts) {
      console.error("key is not defined on this file");
      return;
    }

    if (pathParts.length > 0) {
      dst = pathParts.join("/");
      const dirMap = filesMap.get(dst);
      if (!dirMap) {
        console.error("directory path is not exist");
        return;
      }
    } else {
      // set back to empty str because we dont need the dot when moving to root dir
      dst = "";
    }

    const [error, _] = await promiseResult(Move(dragNode.key!.toString(), dst));
    if (error) {
      console.error(error);
      return;
    }
    console.info(`successfully moved ${dragNode.key?.toString()} to ${dst}`);
  }
  await nextTick();
  emit("get-snippet-paths");
};

const nodeProps = ({ option }: { option: TreeOption }) => {
  return {
    onClick() {
      if (option.isLeaf) {
        emit("file-click", option.key);
      }
    },
  };
};

watchEffect(() => {
  if (props.dirNode) {
    dataRef.value = createData(props.dirNode);
  }
});
</script>

<template>
  <n-tree
    :data="dataRef"
    block-line
    expand-on-click
    draggable
    :node-props="nodeProps"
    @drop="handleDrop"
  />
</template>
