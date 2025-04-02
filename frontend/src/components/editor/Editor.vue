<script setup lang="ts">
import { ref, watch, onMounted, nextTick, shallowRef } from "vue";
import { LanguageSupport } from "@codemirror/language";
import CodeMirror from "vue-codemirror6";
import { NSelect, NIcon } from "naive-ui";
import { OpenPanelLeft } from "@vicons/carbon";
import { basicSetup } from "codemirror";
import type { Extension } from "@codemirror/state";
import { UpdateFile } from "../../../wailsjs/go/main/App.js";
import { promiseResult } from "../../utils/promiseResult.js";
import { themeExtension } from "./theme.js";
import { SelectMixedOption } from "naive-ui/es/select/src/interface.js";

const props = defineProps({
  value: String,
  activePath: String,
});

const emit = defineEmits(["toggle-split"]);

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

const availableLanguages = ref<SelectMixedOption[]>([
  { key: "Go", value: "go" },
  { key: "Javascript", value: "javascript" },
]);
const selectedLanguage = ref("go");

// use shallowRef to avoid typescript yelling type inference too deep
const extensions = shallowRef<Extension[]>([basicSetup, themeExtension]);

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
    extensions.value = [basicSetup, langSupport, themeExtension];
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
  <div class="header">
    <n-icon size="22" @click="emit('toggle-split')">
      <OpenPanelLeft />
    </n-icon>
    <span style="display: inline-block; font-size: 14px; font-weight: 500">
      {{ activePath?.split("/").slice(-1).join("") }}
    </span>
  </div>
  <code-mirror :extensions="extensions" v-model="input" />
  <div class="lang">
    <n-select
      :bordered="false"
      :placeholder="selectedLanguage"
      v-model:value="selectedLanguage"
      :options="availableLanguages"
      label-field="key"
      value-field="value"
    />
  </div>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  gap: 20px;
  width: 100%;
  min-height: 22px;
  padding: 16px;
  border-bottom: 1px solid var(--sumiink-9);
}
.lang {
  position: absolute;
  border-top: 1px solid var(--sumiink-9);
  width: 100%;
  bottom: 0;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}
.lang .n-select {
  width: 130px;
  margin-left: auto;
}

.lang :deep(.n-base-selection-input) {
  left: unset !important;
  display: flex;
  justify-content: end;
  padding: 0 32px 0 12px;
}

/* .lang :deep(.n-base-selection .n-base-selection__border, ) {
  display: none !important;
} */
</style>
