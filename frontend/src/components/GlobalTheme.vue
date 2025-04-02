<script setup lang="ts">
import { darkTheme, NConfigProvider } from "naive-ui";
import type { GlobalThemeOverrides } from "naive-ui";

const themeOverrides: GlobalThemeOverrides = {
  Button: {
    textColor: "red",
  },
  Select: {
    peers: {
      InternalSelection: {
        textColor: "var(--fuji-5)", // Text color inside the input
        color: "var(--sumiink-6)", // Background color of the select input
        placeholderColor: "var(--fuji-5)", // Placeholder text color
        border: "0px",
      },
      InternalSelectMenu: {
        color: "var(--sumiink-5)", // Background color of the dropdown menu
        optionTextColor: "var(--fuji-5)",
        optionTextColorActive: "var(--fuji-4)",
        common: {
          textColor1: "var(--fuji-5)",
          hoverColor: "var(--sumiink-6)",
        },
      },
    },
  },
};
</script>

<template>
  <n-config-provider
    :theme="darkTheme"
    style="height: 100vh"
    :theme-overrides="themeOverrides"
  >
    <slot></slot>
  </n-config-provider>
</template>

<style>
/*split*/
.n-split__resize-trigger {
  background-color: var(--sumiink-9) !important;
}

.n-tree.n-tree--block-line .n-tree-node:not(.n-tree-node--disabled):hover {
  background: var(--sumiink-6) !important;
}

.n-tree.n-tree--block-line
  .n-tree-node:not(.n-tree-node--disabled).n-tree-node--selected {
  background: var(--sumiink-6) !important;
}

.n-split-pane-2 {
  position: relative;
  background-color: var(--sumiink-6);
}

/**tree */
.n-tree-node-switcher {
  display: none !important;
}
.n-tree-node-content__text {
  color: var(--fuji-5) !important;
}
</style>
