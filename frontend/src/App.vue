<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { NButton } from "naive-ui";
import {
  CreateFile,
  GetSnippetPaths,
  ReadFile,
} from "../wailsjs/go/main/App.js";
import Editor from "./components/Editor.vue";
import { promiseResult } from "./utils/promiseResult.js";

const snippets = ref<string[]>([]);
const snippet = ref<string>("");
const activePath = ref<string>("");
const newfileInput = ref<string>("");

const getSnippetPaths = async (): Promise<Array<string>> => {
  const [error, paths] = await promiseResult(GetSnippetPaths());
  if (error) {
    console.error("error getting paths:", error.message);
    return [];
  }
  snippets.value = paths;
  return paths;
};

const readSnippet = async (path: string): Promise<string> => {
  const [error, value] = await promiseResult(ReadFile(path));
  if (error) {
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
  <ul>
    <li
      v-for="(path, idx) in snippets"
      :key="idx"
      @click="handleClickSnippet(path)"
    >
      {{ path }}
    </li>
  </ul>
  <input type="text" placeholder="new file name" v-model="newfileInput" />
  <n-button @click="handleCreateSnippet()">create</n-button>
  <Editor :active-path="activePath" :value="snippet" />
</template>
