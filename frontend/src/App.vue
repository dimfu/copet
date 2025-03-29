<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { NButton } from "naive-ui";
import {
  CreateFile,
  GetSnippetPaths,
  ReadFile,
} from "../wailsjs/go/main/App.js";
import Editor from "./components/Editor.vue";
import FileTreeNode from "./components/FileTreeNode.vue";
import { promiseResult } from "./utils/promiseResult.js";
import { main } from "../wailsjs/go/models.js";

const snippets = ref<main.dirNode>({ files: [], subdirs: {} });
const snippet = ref<string>("");
const activePath = ref<string>("");
const newfileInput = ref<string>("");

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

onMounted(getSnippetPaths);
</script>

<template>
  <FileTreeNode :dir-node="snippets" @file-click="handleClickSnippet" />
  <input type="text" placeholder="new file name" v-model="newfileInput" />
  <n-button @click="handleCreateSnippet()">create</n-button>
  <Editor :active-path="activePath" :value="snippet" />
</template>
