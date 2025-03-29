<script lang="ts" setup>
import { main } from "../../wailsjs/go/models";

const props = defineProps<{ dirNode: main.dirNode }>();

const emit = defineEmits(["file-click"]);
</script>

<template>
  <ul>
    <li
      v-for="(subdir, name) in props.dirNode!.subdirs"
      :key="name"
      @file-click="emit('file-click', $event)"
    >
      📂 {{ name }}
      <FileTreeNode
        @file-click="emit('file-click', $event)"
        :dir-node="subdir"
      />
    </li>
    <li
      v-for="file in props.dirNode!.files"
      :key="file"
      @click="emit('file-click', file)"
    >
      📄 {{ file }}
    </li>
  </ul>
</template>
