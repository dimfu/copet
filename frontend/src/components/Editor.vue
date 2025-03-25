<script setup lang="ts">
import { ref, watch, onMounted, nextTick, shallowRef } from "vue";
import { LanguageSupport } from "@codemirror/language";
import CodeMirror from "vue-codemirror6";
import { basicSetup } from "codemirror";
import type { Extension } from "@codemirror/state";

const value = ref<string>(`fmt.Println("Hello world")`);
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
  <code-mirror :extensions="extensions" v-model="value" />
</template>
