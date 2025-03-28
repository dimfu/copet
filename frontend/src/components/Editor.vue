<script setup lang="ts">
import { ref, watch, onMounted, nextTick, shallowRef, computed } from "vue";
import { LanguageSupport } from "@codemirror/language";
import CodeMirror from "vue-codemirror6";
import { basicSetup } from "codemirror";
import type { Extension } from "@codemirror/state";
import { UpdateFile } from "../../wailsjs/go/main/App.js";
import { promiseResult } from "../utils/promiseResult.js";

const props = defineProps({
  value: String,
  activePath: String,
});

// debounced value update
const timeout = ref<number | null>(null);
const input = ref<string>(props.value ?? "");
const debouncedInput = ref<string>("");
watch(input, (newValue: string) => {
  if (timeout.value) clearTimeout(timeout.value);
  timeout.value = setTimeout(() => {
    debouncedInput.value = newValue;
    if (props.activePath!.length != 0) {
      (async () => {
        const [error] = await promiseResult(
          UpdateFile(props.activePath!, debouncedInput.value)
        );
        if (error) {
          console.error("error while saving file:", error);
        }
      })();
    }
  }, 300);
});

const availableLanguages = ref(["go", "javascript"]);
const selectedLanguage = ref("go");

// use shallowRef to avoid typescript yelling type inference too deep
const extensions = shallowRef<Extension[]>([basicSetup]);

const loadLang = async (lang: string): Promise<LanguageSupport | null> => {
  try {
    switch (lang) {
      case "go":
        return await import("@codemirror/lang-go").then((m) => m.go());
      case "javascript":
        return await import("@codemirror/lang-javascript").then((m) =>
          m.javascript()
        );
      default:
        throw new Error(`Language ${lang} is not supported`);
    }
  } catch (error) {
    console.error("Error loading language:", error);
    return null;
  }
};

const updateExtensions = async (lang: string) => {
  const langSupport = await loadLang(lang);
  if (langSupport) {
    extensions.value = [basicSetup, langSupport];
    await nextTick();
  }
};

onMounted(() => {
  updateExtensions(selectedLanguage.value);
});

watch(selectedLanguage, async (newLang) => {
  updateExtensions(newLang);
});

// when changing file it suppose to update the editor
watch(
  () => props.value,
  (newValue) => {
    input.value = newValue ?? "";
  }
);
</script>

<template>
  <select v-model="selectedLanguage">
    <option
      v-for="language in availableLanguages"
      :key="language"
      :value="language"
    >
      {{ language }}
    </option>
  </select>
  <code-mirror :extensions="extensions" v-model="input" />
</template>
