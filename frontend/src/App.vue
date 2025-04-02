<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { NButton, NSplit } from "naive-ui";
import {
  CreateFile,
  GetSnippetPaths,
  ReadFile,
} from "../wailsjs/go/main/App.js";
import Editor from "./components/editor/Editor.vue";
import FileTreeNode from "./components/FileTreeNode.vue";
import GlobalTheme from "./components/GlobalTheme.vue";
import { promiseResult } from "./utils/promiseResult.js";
import { main } from "../wailsjs/go/models.js";

const snippets = ref<main.dirNode>({ files: [], subdirs: {} });
const snippet = ref<string>("");
const activePath = ref<string>("");
const newfileInput = ref<string>("");

const DEFAULT_SPLIT = 0.15;
// its a hack, dw about it without 0.01 the split aint going anywhere at the first place
const split = ref<number>(DEFAULT_SPLIT + 0.01);

const getSnippetPaths = async () => {
  const [error, node] = await promiseResult(GetSnippetPaths());
  if (error) {
    console.error("error getting paths:", error.message);
    return [];
  }
  // serialize transformed object to class dirnode
  snippets.value = Object.assign(new main.dirNode(), node);
};

const readSnippet = async (path: string): Promise<string> => {
  const [error, value] = await promiseResult(ReadFile(path));
  if (error) {
    console.error(error);
    return "";
  }
  return value;
};

const handleClickSnippet = async (path: string) => {
  const [error, value] = await promiseResult(readSnippet(path));
  if (error) {
    console.error("error reading snippet:", error);
    return;
  }
  snippet.value = value;
  activePath.value = path;
};

// temporary file creation
const handleCreateSnippet = async () => {
  if (newfileInput.value.length == 0) {
    console.error("file name should have atleast 1 character");
    return;
  }
  let [error, path] = await promiseResult(CreateFile(newfileInput.value));
  if (error) {
    console.error("error reading snippet:", error);
    return;
  }

  if (!path || path?.length == 0) {
    console.error("expecting path to be a valid path with set length");
    return;
  }

  await getSnippetPaths();
  await readSnippet(path);
};

const handleDragMove = (event: Event) => {
  const mEvent = event as MouseEvent;
  const pane = document.querySelector(".n-split-pane-1") as HTMLElement;
  if (pane) {
    if (mEvent.clientX <= pane.clientWidth / 2) {
      split.value = 0.0;
    }
  }
};

const handleUpdateSize = (val: string | number) => {
  val = Number(val);
  if (isNaN(val)) return;

  if (
    (split.value != 0 && val > DEFAULT_SPLIT / 2 - 0.0001) ||
    (split.value == 0 && val > DEFAULT_SPLIT)
  ) {
    split.value = val;
  }
};

onMounted(getSnippetPaths);
</script>

<template>
  <GlobalTheme>
    <n-split
      v-model:size="split"
      :on-drag-move="handleDragMove"
      :on-update:size="handleUpdateSize"
      direction="horizontal"
      :max="1"
      :min="0.15"
      :resize-trigger-size="1"
    >
      <template #1>
        <div class="header"></div>
        <div>
          <input
            type="text"
            placeholder="new file name"
            v-model="newfileInput"
          />
          <n-button @click="handleCreateSnippet()">create</n-button>
          <FileTreeNode
            :dir-node="snippets"
            @file-click="handleClickSnippet"
            @get-snippet-paths="getSnippetPaths"
          />
        </div>
      </template>
      <template #2>
        <Editor :active-path="activePath" :value="snippet" />
      </template>
    </n-split>
  </GlobalTheme>
</template>

<style scoped>
.header {
  min-height: 22px;
  padding: 16px;
  border-bottom: 1px solid var(--sumiink-9);
}
</style>
