import { TagStyle } from "@codemirror/language";
import { createTheme } from "@uiw/codemirror-themes";
import { tags as t } from "@lezer/highlight";
import { colors } from "./colors";

export const tagStyles: TagStyle[] = [
  { tag: t.variableName, color: colors["editor.fg"] },
  { tag: t.content, backgroundColor: colors["editor.bg"] },
  { tag: t.meta, backgroundColor: colors["editor.line"] },

  { tag: t.tagName, color: colors["syntax.tag"] },
  { tag: t.function(t.variableName), color: colors["syntax.func"] },
  { tag: t.className, color: colors["syntax.entity"] },
  { tag: t.string, color: colors["syntax.string"] },
  { tag: t.regexp, color: colors["syntax.regexp"] },
  { tag: t.monospace, color: colors["syntax.markup"] },
  { tag: t.keyword, color: colors["syntax.keyword"] },
  { tag: t.special(t.variableName), color: colors["syntax.special"] },
  { tag: t.comment, color: colors["syntax.comment"] },
  { tag: t.constant(t.variableName), color: colors["syntax.constant"] },
  { tag: t.operator, color: colors["syntax.operator"] },

  { tag: t.inserted, color: colors["vcs.added"] },
  { tag: t.changed, color: colors["vcs.modified"] },
  { tag: t.deleted, color: colors["vcs.removed"] },
];

export const themeExtension = createTheme({
  theme: "dark",
  settings: {
    background: colors["editor.bg1"],
    backgroundImage: "",
    foreground: colors["editor.fg"],
    gutterBackground: "inherit",
    gutterForeground: "inherit",
  },
  styles: tagStyles,
});
