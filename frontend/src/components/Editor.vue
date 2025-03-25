<script lang="ts">
import { ref, defineComponent, watch, onMounted } from "vue";
import { LanguageSupport } from "@codemirror/language";

import CodeMirror from "vue-codemirror6";
import { basicSetup } from "codemirror";

export default defineComponent({
  components: {
    CodeMirror,
  },
  setup() {
    const value = ref<string>(`fmt.Println("Hello world")`);
    const availableLanguages = ref(["go", "javascript"]);
    const selectedLanguage = ref("go");
    let extensions = ref<any[]>([]);

    const loadLang = async (
      lang: string
    ): Promise<LanguageSupport | undefined> => {
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
      }
    };

    onMounted(async () => {
      const module = await loadLang(selectedLanguage.value);
      if (module) {
        extensions.value = [basicSetup, module];
      }
    });

    watch(selectedLanguage, async (newLang) => {
      const module = await loadLang(newLang);
      if (module) {
        console.log("succesfully changed to ", newLang);
        extensions.value = [basicSetup, module];
      }
    });

    return { availableLanguages, selectedLanguage, extensions, value };
  },
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
