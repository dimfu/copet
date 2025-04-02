<script lang="ts" setup>
import type { TreeDropInfo, TreeOption } from "naive-ui";
import { NTree, NIcon } from "naive-ui";
import { Folder, FolderOpen } from "@vicons/carbon";
import { main } from "../../wailsjs/go/models";
import { h, nextTick, onMounted, ref, watchEffect } from "vue";
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
      label: fileName,
      key: file,
      isLeaf: true,
      children: [],
    };
  });

  const subdirs: TreeOptionWithDir[] = Object.entries(node.subdirs || {}).map(
    ([name, subdir]) => {
      const currentPath = k != undefined ? `${k}/${name}` : name;

      // initialize map to avoid skipping empty folder from being listed in map
      if (!filesMap.has(currentPath)) {
        filesMap.set(currentPath, []);
      }

      return {
        directory: name,
        label: name,
        key: currentPath,
        isLeaf: false,
        prefix: () =>
          h(NIcon, null, {
            default: () => h(Folder),
          }),
        children: createData(subdir, name),
      };
    }
  );

  return [...subdirs, ...files];
};

const findSiblingsAndIndex = (
  node: TreeOptionWithDir,
  nodes?: TreeOptionWithDir[]
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
  try {
    if (dropPosition == "inside" && !node.isLeaf) {
      const [error, _] = await promiseResult(
        Move(dragNode.key!.toString(), node.key!.toString())
      );
      if (error) {
        console.error(error);
        return;
      }
    } else {
      const [nodeSiblings, nodeIndex] = findSiblingsAndIndex(
        node,
        dataRef.value
      );
      if (nodeSiblings === null || nodeIndex === null) {
        return;
      }

      let dst: string = ".";
      const pathParts = node.key?.toString().split("/").slice(0, -1);
      if (!pathParts) {
        console.error("key is not defined on this file");
        return;
      }

      if (pathParts.length > 0) {
        dst = pathParts.join("/");
        if (dst == dragNode.directory) {
          return;
        }
        const files = filesMap.get(dst);
        if (!files) {
          console.error("directory path does not exist");
          return;
        }
        // prevent file/folder from moving to their current directory
        const curr = filesMap.get(`${dst}/${dragNode.label}`);
        if (curr) {
          return;
        }
      } else {
        dst = "";
      }

      const [error, _] = await promiseResult(
        Move(dragNode.key!.toString(), dst)
      );
      if (error) {
        console.error(error);
        return;
      }
      console.info(`successfully moved ${dragNode.key?.toString()} to ${dst}`);
    }
  } finally {
    await nextTick();
    emit("get-snippet-paths");
    filesMap.clear();
  }
};

const updatePrefixWithExpaned = (
  _keys: Array<string | number>,
  _option: Array<TreeOption | null>,
  meta: {
    node: TreeOption | null;
    action: "expand" | "collapse" | "filter";
  }
) => {
  if (!meta.node) return;
  switch (meta.action) {
    case "expand":
      meta.node.prefix = () =>
        h(NIcon, null, {
          default: () => h(FolderOpen),
        });
      break;
    case "collapse":
      meta.node.prefix = () =>
        h(NIcon, null, {
          default: () => h(Folder),
        });
      break;
  }
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
    show-line
    draggable
    :render-switcher-icon="() => null"
    :on-update:expanded-keys="updatePrefixWithExpaned"
    :node-props="nodeProps"
    @drop="handleDrop"
  />
</template>
