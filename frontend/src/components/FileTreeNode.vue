<script lang="ts" setup>
import type { TreeOption } from "naive-ui";
import { NTree } from "naive-ui";
import { main } from "../../wailsjs/go/models";

const props = defineProps<{ dirNode: main.dirNode }>();
const emit = defineEmits(["file-click"]);

const createData = (node: main.dirNode, k: string): Array<TreeOption> => {
  if (!node) return [];

  const files = node.files.map((file) => ({
    label: file,
    key: file,
    isLeaf: true,
    children: [],
  }));

  const subdirs = Object.entries(node.subdirs || {}).map(([name, subdir]) => ({
    label: name,
    key: `${k}_${name}`,
    children: createData(subdir, name),
  }));

  return [...subdirs, ...files];
};

const nodeProps = ({ option }: { option: TreeOption }) => {
  return {
    onClick() {
      if (option.isLeaf) {
        emit("file-click", option.label);
      }
    },
  };
};
</script>

<template>
  <n-tree
    block-line
    :data="createData(props.dirNode, 'root')"
    :node-props="nodeProps"
    expand-on-click
  />
</template>
